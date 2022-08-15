// Code generated by go-swagger; DO NOT EDIT.

package network

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"visjs/gen/models"
)

// NwOKCode is the HTTP code returned for type NwOK
const NwOKCode int = 200

/*NwOK A successful response.

swagger:response nwOK
*/
type NwOK struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewNwOK creates NwOK with default headers values
func NewNwOK() *NwOK {

	return &NwOK{}
}

// WithPayload adds the payload to the nw o k response
func (o *NwOK) WithPayload(payload interface{}) *NwOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the nw o k response
func (o *NwOK) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *NwOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*NwDefault An unexpected error response.

swagger:response nwDefault
*/
type NwDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.RPCStatus `json:"body,omitempty"`
}

// NewNwDefault creates NwDefault with default headers values
func NewNwDefault(code int) *NwDefault {
	if code <= 0 {
		code = 500
	}

	return &NwDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the nw default response
func (o *NwDefault) WithStatusCode(code int) *NwDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the nw default response
func (o *NwDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the nw default response
func (o *NwDefault) WithPayload(payload *models.RPCStatus) *NwDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the nw default response
func (o *NwDefault) SetPayload(payload *models.RPCStatus) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *NwDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}