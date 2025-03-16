// Code generated by go-swagger; DO NOT EDIT.

package config

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

// NewFindConfigPrometheusAddressParams creates a new FindConfigPrometheusAddressParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigPrometheusAddressParams() *FindConfigPrometheusAddressParams {
	return &FindConfigPrometheusAddressParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigPrometheusAddressParamsWithTimeout creates a new FindConfigPrometheusAddressParams object
// with the ability to set a timeout on a request.
func NewFindConfigPrometheusAddressParamsWithTimeout(timeout time.Duration) *FindConfigPrometheusAddressParams {
	return &FindConfigPrometheusAddressParams{
		timeout: timeout,
	}
}

// NewFindConfigPrometheusAddressParamsWithContext creates a new FindConfigPrometheusAddressParams object
// with the ability to set a context for a request.
func NewFindConfigPrometheusAddressParamsWithContext(ctx context.Context) *FindConfigPrometheusAddressParams {
	return &FindConfigPrometheusAddressParams{
		Context: ctx,
	}
}

// NewFindConfigPrometheusAddressParamsWithHTTPClient creates a new FindConfigPrometheusAddressParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigPrometheusAddressParamsWithHTTPClient(client *http.Client) *FindConfigPrometheusAddressParams {
	return &FindConfigPrometheusAddressParams{
		HTTPClient: client,
	}
}

/*
FindConfigPrometheusAddressParams contains all the parameters to send to the API endpoint

	for the find config prometheus address operation.

	Typically these are written to a http.Request.
*/
type FindConfigPrometheusAddressParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config prometheus address params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigPrometheusAddressParams) WithDefaults() *FindConfigPrometheusAddressParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config prometheus address params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigPrometheusAddressParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) WithTimeout(timeout time.Duration) *FindConfigPrometheusAddressParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) WithContext(ctx context.Context) *FindConfigPrometheusAddressParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) WithHTTPClient(client *http.Client) *FindConfigPrometheusAddressParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config prometheus address params
func (o *FindConfigPrometheusAddressParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigPrometheusAddressParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
