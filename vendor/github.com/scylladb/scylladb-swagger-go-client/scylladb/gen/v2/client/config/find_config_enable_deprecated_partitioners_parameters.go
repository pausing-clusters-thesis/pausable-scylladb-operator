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

// NewFindConfigEnableDeprecatedPartitionersParams creates a new FindConfigEnableDeprecatedPartitionersParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigEnableDeprecatedPartitionersParams() *FindConfigEnableDeprecatedPartitionersParams {
	return &FindConfigEnableDeprecatedPartitionersParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigEnableDeprecatedPartitionersParamsWithTimeout creates a new FindConfigEnableDeprecatedPartitionersParams object
// with the ability to set a timeout on a request.
func NewFindConfigEnableDeprecatedPartitionersParamsWithTimeout(timeout time.Duration) *FindConfigEnableDeprecatedPartitionersParams {
	return &FindConfigEnableDeprecatedPartitionersParams{
		timeout: timeout,
	}
}

// NewFindConfigEnableDeprecatedPartitionersParamsWithContext creates a new FindConfigEnableDeprecatedPartitionersParams object
// with the ability to set a context for a request.
func NewFindConfigEnableDeprecatedPartitionersParamsWithContext(ctx context.Context) *FindConfigEnableDeprecatedPartitionersParams {
	return &FindConfigEnableDeprecatedPartitionersParams{
		Context: ctx,
	}
}

// NewFindConfigEnableDeprecatedPartitionersParamsWithHTTPClient creates a new FindConfigEnableDeprecatedPartitionersParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigEnableDeprecatedPartitionersParamsWithHTTPClient(client *http.Client) *FindConfigEnableDeprecatedPartitionersParams {
	return &FindConfigEnableDeprecatedPartitionersParams{
		HTTPClient: client,
	}
}

/*
FindConfigEnableDeprecatedPartitionersParams contains all the parameters to send to the API endpoint

	for the find config enable deprecated partitioners operation.

	Typically these are written to a http.Request.
*/
type FindConfigEnableDeprecatedPartitionersParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config enable deprecated partitioners params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigEnableDeprecatedPartitionersParams) WithDefaults() *FindConfigEnableDeprecatedPartitionersParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config enable deprecated partitioners params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigEnableDeprecatedPartitionersParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) WithTimeout(timeout time.Duration) *FindConfigEnableDeprecatedPartitionersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) WithContext(ctx context.Context) *FindConfigEnableDeprecatedPartitionersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) WithHTTPClient(client *http.Client) *FindConfigEnableDeprecatedPartitionersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config enable deprecated partitioners params
func (o *FindConfigEnableDeprecatedPartitionersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigEnableDeprecatedPartitionersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
