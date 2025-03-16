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

// NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams creates a new ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams() *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	return &ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithTimeout creates a new ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams object
// with the ability to set a timeout on a request.
func NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithTimeout(timeout time.Duration) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	return &ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams{
		timeout: timeout,
	}
}

// NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithContext creates a new ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams object
// with the ability to set a context for a request.
func NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithContext(ctx context.Context) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	return &ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams{
		Context: ctx,
	}
}

// NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithHTTPClient creates a new ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParamsWithHTTPClient(client *http.Client) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	return &ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams{
		HTTPClient: client,
	}
}

/*
ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams contains all the parameters to send to the API endpoint

	for the column family metrics total disk space used by name get operation.

	Typically these are written to a http.Request.
*/
type ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams struct {

	/* Name.

	   The column family name in keyspace:name format
	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the column family metrics total disk space used by name get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WithDefaults() *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the column family metrics total disk space used by name get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WithTimeout(timeout time.Duration) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WithContext(ctx context.Context) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WithHTTPClient(client *http.Client) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WithName(name string) *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the column family metrics total disk space used by name get params
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *ColumnFamilyMetricsTotalDiskSpaceUsedByNameGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
