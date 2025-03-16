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

// NewCacheServiceMetricsKeyHitsGetParams creates a new CacheServiceMetricsKeyHitsGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCacheServiceMetricsKeyHitsGetParams() *CacheServiceMetricsKeyHitsGetParams {
	return &CacheServiceMetricsKeyHitsGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCacheServiceMetricsKeyHitsGetParamsWithTimeout creates a new CacheServiceMetricsKeyHitsGetParams object
// with the ability to set a timeout on a request.
func NewCacheServiceMetricsKeyHitsGetParamsWithTimeout(timeout time.Duration) *CacheServiceMetricsKeyHitsGetParams {
	return &CacheServiceMetricsKeyHitsGetParams{
		timeout: timeout,
	}
}

// NewCacheServiceMetricsKeyHitsGetParamsWithContext creates a new CacheServiceMetricsKeyHitsGetParams object
// with the ability to set a context for a request.
func NewCacheServiceMetricsKeyHitsGetParamsWithContext(ctx context.Context) *CacheServiceMetricsKeyHitsGetParams {
	return &CacheServiceMetricsKeyHitsGetParams{
		Context: ctx,
	}
}

// NewCacheServiceMetricsKeyHitsGetParamsWithHTTPClient creates a new CacheServiceMetricsKeyHitsGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewCacheServiceMetricsKeyHitsGetParamsWithHTTPClient(client *http.Client) *CacheServiceMetricsKeyHitsGetParams {
	return &CacheServiceMetricsKeyHitsGetParams{
		HTTPClient: client,
	}
}

/*
CacheServiceMetricsKeyHitsGetParams contains all the parameters to send to the API endpoint

	for the cache service metrics key hits get operation.

	Typically these are written to a http.Request.
*/
type CacheServiceMetricsKeyHitsGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the cache service metrics key hits get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CacheServiceMetricsKeyHitsGetParams) WithDefaults() *CacheServiceMetricsKeyHitsGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the cache service metrics key hits get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CacheServiceMetricsKeyHitsGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) WithTimeout(timeout time.Duration) *CacheServiceMetricsKeyHitsGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) WithContext(ctx context.Context) *CacheServiceMetricsKeyHitsGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) WithHTTPClient(client *http.Client) *CacheServiceMetricsKeyHitsGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the cache service metrics key hits get params
func (o *CacheServiceMetricsKeyHitsGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *CacheServiceMetricsKeyHitsGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
