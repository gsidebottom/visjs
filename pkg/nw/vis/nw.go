package vis

type Node struct {
	Id    int    `json:"id"`
	Label string `json:"label,omitempty"`
}

type Nodes []*Node

type Link struct {
	Id     int    `json:"id,omitempty"`
	Label  string `json:"label,omitempty"`
	Length int    `json:"length,omitempty"` // spring length for visjs physics
	Width  int    `json:"width,omitempty"`
	From   int    `json:"from"`
	To     int    `json:"to"`
}

type Links []*Link

type Network struct {
	Nodes `json:"nodes,omitempty"`
	Links `json:"edges,omitempty"`
}
