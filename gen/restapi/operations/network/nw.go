// Code generated by go-swagger; DO NOT EDIT.

package network

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// NwHandlerFunc turns a function with the right signature into a nw handler
type NwHandlerFunc func(NwParams) middleware.Responder

// Handle executing the request and returning a response
func (fn NwHandlerFunc) Handle(params NwParams) middleware.Responder {
	return fn(params)
}

// NwHandler interface for that can handle valid nw params
type NwHandler interface {
	Handle(NwParams) middleware.Responder
}

// NewNw creates a new http.Handler for the nw operation
func NewNw(ctx *middleware.Context, handler NwHandler) *Nw {
	return &Nw{Context: ctx, Handler: handler}
}

/* Nw swagger:route GET /api/v1/nw/{nwId} Network nw

Get Network

Get Network

*/
type Nw struct {
	Context *middleware.Context
	Handler NwHandler
}

func (o *Nw) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewNwParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
