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
	"github.com/go-openapi/swag"
)

// NewColumnFamilyCompactionByNamePostParams creates a new ColumnFamilyCompactionByNamePostParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewColumnFamilyCompactionByNamePostParams() *ColumnFamilyCompactionByNamePostParams {
	return &ColumnFamilyCompactionByNamePostParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewColumnFamilyCompactionByNamePostParamsWithTimeout creates a new ColumnFamilyCompactionByNamePostParams object
// with the ability to set a timeout on a request.
func NewColumnFamilyCompactionByNamePostParamsWithTimeout(timeout time.Duration) *ColumnFamilyCompactionByNamePostParams {
	return &ColumnFamilyCompactionByNamePostParams{
		timeout: timeout,
	}
}

// NewColumnFamilyCompactionByNamePostParamsWithContext creates a new ColumnFamilyCompactionByNamePostParams object
// with the ability to set a context for a request.
func NewColumnFamilyCompactionByNamePostParamsWithContext(ctx context.Context) *ColumnFamilyCompactionByNamePostParams {
	return &ColumnFamilyCompactionByNamePostParams{
		Context: ctx,
	}
}

// NewColumnFamilyCompactionByNamePostParamsWithHTTPClient creates a new ColumnFamilyCompactionByNamePostParams object
// with the ability to set a custom HTTPClient for a request.
func NewColumnFamilyCompactionByNamePostParamsWithHTTPClient(client *http.Client) *ColumnFamilyCompactionByNamePostParams {
	return &ColumnFamilyCompactionByNamePostParams{
		HTTPClient: client,
	}
}

/*
ColumnFamilyCompactionByNamePostParams contains all the parameters to send to the API endpoint

	for the column family compaction by name post operation.

	Typically these are written to a http.Request.
*/
type ColumnFamilyCompactionByNamePostParams struct {

	/* Maximum.

	   The maximum number of sstables in queue before compaction kicks off

	   Format: int32
	*/
	Maximum int32

	/* Minimum.

	   The minimum number of sstables in queue before compaction kicks off

	   Format: int32
	*/
	Minimum int32

	/* Name.

	   The column family name in keyspace:name format
	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the column family compaction by name post params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ColumnFamilyCompactionByNamePostParams) WithDefaults() *ColumnFamilyCompactionByNamePostParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the column family compaction by name post params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ColumnFamilyCompactionByNamePostParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithTimeout(timeout time.Duration) *ColumnFamilyCompactionByNamePostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithContext(ctx context.Context) *ColumnFamilyCompactionByNamePostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithHTTPClient(client *http.Client) *ColumnFamilyCompactionByNamePostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithMaximum adds the maximum to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithMaximum(maximum int32) *ColumnFamilyCompactionByNamePostParams {
	o.SetMaximum(maximum)
	return o
}

// SetMaximum adds the maximum to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetMaximum(maximum int32) {
	o.Maximum = maximum
}

// WithMinimum adds the minimum to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithMinimum(minimum int32) *ColumnFamilyCompactionByNamePostParams {
	o.SetMinimum(minimum)
	return o
}

// SetMinimum adds the minimum to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetMinimum(minimum int32) {
	o.Minimum = minimum
}

// WithName adds the name to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) WithName(name string) *ColumnFamilyCompactionByNamePostParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the column family compaction by name post params
func (o *ColumnFamilyCompactionByNamePostParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *ColumnFamilyCompactionByNamePostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param maximum
	qrMaximum := o.Maximum
	qMaximum := swag.FormatInt32(qrMaximum)
	if qMaximum != "" {

		if err := r.SetQueryParam("maximum", qMaximum); err != nil {
			return err
		}
	}

	// query param minimum
	qrMinimum := o.Minimum
	qMinimum := swag.FormatInt32(qrMinimum)
	if qMinimum != "" {

		if err := r.SetQueryParam("minimum", qMinimum); err != nil {
			return err
		}
	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
