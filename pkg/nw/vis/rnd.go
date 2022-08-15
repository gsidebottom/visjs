package vis

import (
	"fmt"
	"github.com/fogleman/delaunay"
	"github.com/fogleman/poissondisc"
	"math"
	"math/rand"
	"sort"
)

func NewRandNetwork(numNodes int) (nw Network) {
	nw.Nodes = make(Nodes, numNodes)
	connectionCount := make([]int, numNodes)
	for i := range nw.Nodes {
		nw.Nodes[i] = &Node{
			Id:    i,
			Label: fmt.Sprintf("%d", i),
		}
		switch {
		case i == 0:
			from, to := i, 0
			nw.Links = append(nw.Links, &Link{
				From: from,
				To:   to,
			})
			connectionCount[from]++
			connectionCount[to]++
		case i > 0:
			conn := len(nw.Links) * 2
			r := rand.Intn(conn)
			var cum int
			var j int
			for j < len(connectionCount) && cum < r {
				cum += connectionCount[j]
				j++
			}
			from, to := i, j
			nw.Links = append(nw.Links, &Link{
				From: from,
				To:   to,
			})
			connectionCount[from]++
			connectionCount[to]++
		}
	}

	return
}

type AZNodeNames struct{ A, Z string }
type LinkData struct {
	AM, ZM uint64
	BW     uint64
}
type Lnks map[AZNodeNames]LinkData
type NdCoords map[string][2]float64
type NdsOrLnks []string
type Sites map[string]NdsOrLnks

// network with random planar IP layer, nSites sites with an average of nNodesPerSite.  The members of each site
// form an Anycast group
// siteSpread is factor to spread out sites with, about 3.5 ma
func NewRandPlanarNetwork(nSites, nNodesPerSite int, minSiteDist float64, springFactor int, bwOpts []uint64) (nw *Network) {
	generatePoints := func(n int) ([]delaunay.Point, float64) {
		// r := math.Sqrt(float64(n) * 1.618)
		// find the smallest r that yields a grid that can fit n points
		for r := 1.0; ; r++ {
			points := poissondisc.Sample(-r, -r, r, r, 1, 32, nil)
			if len(points) < n {
				continue
			}
			//rand.Shuffle(len(points), func(i, j int) {
			// points[i], points[j] = points[j], points[i]
			//})
			sort.Slice(points, func(i, j int) bool {
				p1 := points[i]
				p2 := points[j]
				d1 := math.Hypot(p1.X, p1.Y)
				d2 := math.Hypot(p2.X, p2.Y)
				return d1 < d2
			})
			points = points[:n]
			result := make([]delaunay.Point, len(points))
			for i, p := range points {
				result[i].X = p.X
				result[i].Y = p.Y
			}
			return result, r
		}
	}
	nextHalfEdge := func(e int) int {
		if e%3 == 2 {
			return e - 2
		}
		return e + 1
	}

	// generate points for nodes in each site
	ptsBySite := make([][]delaunay.Point, nSites)
	aveSiteRadius := 0.0
	for i := range ptsBySite {
		var sr float64

		ptsBySite[i], sr = generatePoints(int(rand.NormFloat64() + float64(nNodesPerSite)))
		aveSiteRadius += sr
	}
	aveSiteRadius /= float64(nSites)

	// generate points corresponding to centers of sites
	siteCenterPts, _ := generatePoints(nSites)
	for i, siteCenterPt := range siteCenterPts {
		siteCenterPt.X *= aveSiteRadius*2 + (minSiteDist - 1) // min site distance is 1 already, add zero if we want 1
		siteCenterPt.Y *= aveSiteRadius*2 + (minSiteDist - 1)
		siteCenterPts[i] = siteCenterPt // by value, need to put value back
	}

	// translate each site to it's center
	for i, sitePts := range ptsBySite {
		siteCenterPt := siteCenterPts[i]

		for j := range sitePts {
			sitePts[j].X += siteCenterPt.X
			sitePts[j].Y += siteCenterPt.Y
		}
	}
	var points []delaunay.Point
	for _, sitePts := range ptsBySite {
		points = append(points, sitePts...)
	}

	ndCoords := NdCoords{}
	coords2Nd := map[[2]float64]string{}

	for sNum, sitePts := range ptsBySite {
		for nNum, p := range sitePts {
			name := fmt.Sprintf("%X_%X", sNum, nNum)
			coords := [2]float64{p.X, p.Y}
			ndCoords[name] = coords
			coords2Nd[coords] = name
		}

	}

	ls := Lnks{}
	addLink := func(aPt, zPt delaunay.Point) {
		dist := func() uint64 {
			dx, dy := aPt.X-zPt.X, aPt.Y-zPt.Y
			d := uint64(math.Sqrt(dx*dx + dy*dy))
			if d == 0 {
				d = 1
			}
			return d
		}
		aCoord, zCoord := [2]float64{aPt.X, aPt.Y}, [2]float64{zPt.X, zPt.Y}

		a, z := coords2Nd[aCoord], coords2Nd[zCoord]
		d := dist()
		bw := bwOpts[rand.Intn(len(bwOpts))]
		ls[AZNodeNames{a, z}] = LinkData{d, d, bw}
	}

	triangulation, _ := delaunay.Triangulate(points)
	ts := triangulation.Triangles
	hs := triangulation.Halfedges
	for i, h := range hs {
		if i > h {
			aPt := points[ts[i]]
			zPt := points[ts[nextHalfEdge(i)]]
			addLink(aPt, zPt)
		}
	}

	nw = NewNW(ndCoords, ls, springFactor)
	return
}

func NewNW(ncs NdCoords, lnks Lnks, springFactor int) (nw *Network) {
	nw = &Network{}
	nw.Nodes = make([]*Node, len(ncs))
	ndLabelToID := map[string]int{}
	var i int
	for ndLabel := range ncs {
		nd := &Node{
			Id:    i + 1,
			Label: ndLabel,
		}
		nw.Nodes[i] = nd
		ndLabelToID[ndLabel] = nd.Id
		i++
	}
	nw.Links = make([]*Link, len(lnks))
	i = 0
	for az, ld := range lnks {
		a, z := ndLabelToID[az.A], ndLabelToID[az.Z]
		lnk := &Link{
			Label:  fmt.Sprintf("%d", ld.AM),
			Length: int(ld.AM*ld.ZM) * springFactor,
			// Length: int(ld.AM) * 20,
			Width: int(ld.BW) / 1e9,
			From:  a,
			To:    z,
		}
		nw.Links[i] = lnk
		i++
	}
	return
}
