// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGossiperEndpointDownGetParams creates a new GossiperEndpointDownGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGossiperEndpointDownGetParams() *GossiperEndpointDownGetParams {
	return &GossiperEndpointDownGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGossiperEndpointDownGetParamsWithTimeout creates a new GossiperEndpointDownGetParams object
// with the ability to set a timeout on a request.
func NewGossiperEndpointDownGetParamsWithTimeout(timeout time.Duration) *GossiperEndpointDownGetParams {
	return &GossiperEndpointDownGetParams{
		timeout: timeout,
	}
}

// NewGossiperEndpointDownGetParamsWithContext creates a new GossiperEndpointDownGetParams object
// with the ability to set a context for a request.
func NewGossiperEndpointDownGetParamsWithContext(ctx context.Context) *GossiperEndpointDownGetParams {
	return &GossiperEndpointDownGetParams{
		Context: ctx,
	}
}

// NewGossiperEndpointDownGetParamsWithHTTPClient creates a new GossiperEndpointDownGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewGossiperEndpointDownGetParamsWithHTTPClient(client *http.Client) *GossiperEndpointDownGetParams {
	return &GossiperEndpointDownGetParams{
		HTTPClient: client,
	}
}

/*
GossiperEndpointDownGetParams contains all the parameters to send to the API endpoint

	for the gossiper endpoint down get operation.

	Typically these are written to a http.Request.
*/
type GossiperEndpointDownGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the gossiper endpoint down get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GossiperEndpointDownGetParams) WithDefaults() *GossiperEndpointDownGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the gossiper endpoint down get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GossiperEndpointDownGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) WithTimeout(timeout time.Duration) *GossiperEndpointDownGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) WithContext(ctx context.Context) *GossiperEndpointDownGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) WithHTTPClient(client *http.Client) *GossiperEndpointDownGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the gossiper endpoint down get params
func (o *GossiperEndpointDownGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GossiperEndpointDownGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
