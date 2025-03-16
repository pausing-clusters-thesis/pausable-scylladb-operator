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

// NewFindConfigCommitlogSyncPeriodInMsParams creates a new FindConfigCommitlogSyncPeriodInMsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigCommitlogSyncPeriodInMsParams() *FindConfigCommitlogSyncPeriodInMsParams {
	return &FindConfigCommitlogSyncPeriodInMsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigCommitlogSyncPeriodInMsParamsWithTimeout creates a new FindConfigCommitlogSyncPeriodInMsParams object
// with the ability to set a timeout on a request.
func NewFindConfigCommitlogSyncPeriodInMsParamsWithTimeout(timeout time.Duration) *FindConfigCommitlogSyncPeriodInMsParams {
	return &FindConfigCommitlogSyncPeriodInMsParams{
		timeout: timeout,
	}
}

// NewFindConfigCommitlogSyncPeriodInMsParamsWithContext creates a new FindConfigCommitlogSyncPeriodInMsParams object
// with the ability to set a context for a request.
func NewFindConfigCommitlogSyncPeriodInMsParamsWithContext(ctx context.Context) *FindConfigCommitlogSyncPeriodInMsParams {
	return &FindConfigCommitlogSyncPeriodInMsParams{
		Context: ctx,
	}
}

// NewFindConfigCommitlogSyncPeriodInMsParamsWithHTTPClient creates a new FindConfigCommitlogSyncPeriodInMsParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigCommitlogSyncPeriodInMsParamsWithHTTPClient(client *http.Client) *FindConfigCommitlogSyncPeriodInMsParams {
	return &FindConfigCommitlogSyncPeriodInMsParams{
		HTTPClient: client,
	}
}

/*
FindConfigCommitlogSyncPeriodInMsParams contains all the parameters to send to the API endpoint

	for the find config commitlog sync period in ms operation.

	Typically these are written to a http.Request.
*/
type FindConfigCommitlogSyncPeriodInMsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config commitlog sync period in ms params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigCommitlogSyncPeriodInMsParams) WithDefaults() *FindConfigCommitlogSyncPeriodInMsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config commitlog sync period in ms params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigCommitlogSyncPeriodInMsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) WithTimeout(timeout time.Duration) *FindConfigCommitlogSyncPeriodInMsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) WithContext(ctx context.Context) *FindConfigCommitlogSyncPeriodInMsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) WithHTTPClient(client *http.Client) *FindConfigCommitlogSyncPeriodInMsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config commitlog sync period in ms params
func (o *FindConfigCommitlogSyncPeriodInMsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigCommitlogSyncPeriodInMsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
