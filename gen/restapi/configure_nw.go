// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"visjs/pkg/nw/vis"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"visjs/gen/restapi/operations"
	"visjs/gen/restapi/operations/network"
)

//go:generate swagger generate server --target ../../gen --name Nw --spec ../../api/openapi.yml --principal interface{}

func configureFlags(api *operations.NwAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }

}

func configureAPI(api *operations.NwAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	errResponse := func(code int, errMsg string) middleware.ResponderFunc {
		b, _ := json.Marshal(map[string]interface{}{
			"error": errMsg,
		})
		return middleware.ResponderFunc(func(writer http.ResponseWriter, producer runtime.Producer) {
			writer.WriteHeader(code)
			writer.Write(b)
		})
	}
	successResponse := func(nw interface{}) middleware.ResponderFunc {
		var b []byte
		var err error
		if b, err = json.Marshal(nw); err != nil {
			return errResponse(http.StatusInternalServerError, fmt.Sprintf("Network Marshal error: %s", err))
		}
		return middleware.ResponderFunc(func(writer http.ResponseWriter, producer runtime.Producer) {
			writer.WriteHeader(http.StatusOK)
			writer.Write(b)
		})
	}

	api.NetworkNwHandler = network.NwHandlerFunc(func(params network.NwParams) middleware.Responder {
		numSites, numNodesPerSite, springFactor := 5, 6, 2
		spl := strings.Split(params.NwID, "_")
		if len(spl) >= 3 {
			var num int
			var err error
			if num, err = strconv.Atoi(spl[0]); err == nil {
				numSites = num
			}
			if num, err = strconv.Atoi(spl[1]); err == nil {
				numNodesPerSite = num
			}
			if num, err = strconv.Atoi(spl[2]); err == nil {
				springFactor = num
			}
		}
		return successResponse(vis.NewRandPlanarNetwork(numSites, numNodesPerSite, 10, springFactor, []uint64{1e9, 4e9, 10e9}))
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return fileServerMiddleware(handler)
}

// fileServerMiddleware serve static files if it's not an API call
func fileServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
		} else {
			http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
		}
	})
}
