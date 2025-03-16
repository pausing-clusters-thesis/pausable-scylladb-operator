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

// NewFindConfigInternodeAuthenticatorParams creates a new FindConfigInternodeAuthenticatorParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewFindConfigInternodeAuthenticatorParams() *FindConfigInternodeAuthenticatorParams {
	return &FindConfigInternodeAuthenticatorParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewFindConfigInternodeAuthenticatorParamsWithTimeout creates a new FindConfigInternodeAuthenticatorParams object
// with the ability to set a timeout on a request.
func NewFindConfigInternodeAuthenticatorParamsWithTimeout(timeout time.Duration) *FindConfigInternodeAuthenticatorParams {
	return &FindConfigInternodeAuthenticatorParams{
		timeout: timeout,
	}
}

// NewFindConfigInternodeAuthenticatorParamsWithContext creates a new FindConfigInternodeAuthenticatorParams object
// with the ability to set a context for a request.
func NewFindConfigInternodeAuthenticatorParamsWithContext(ctx context.Context) *FindConfigInternodeAuthenticatorParams {
	return &FindConfigInternodeAuthenticatorParams{
		Context: ctx,
	}
}

// NewFindConfigInternodeAuthenticatorParamsWithHTTPClient creates a new FindConfigInternodeAuthenticatorParams object
// with the ability to set a custom HTTPClient for a request.
func NewFindConfigInternodeAuthenticatorParamsWithHTTPClient(client *http.Client) *FindConfigInternodeAuthenticatorParams {
	return &FindConfigInternodeAuthenticatorParams{
		HTTPClient: client,
	}
}

/*
FindConfigInternodeAuthenticatorParams contains all the parameters to send to the API endpoint

	for the find config internode authenticator operation.

	Typically these are written to a http.Request.
*/
type FindConfigInternodeAuthenticatorParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the find config internode authenticator params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigInternodeAuthenticatorParams) WithDefaults() *FindConfigInternodeAuthenticatorParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the find config internode authenticator params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *FindConfigInternodeAuthenticatorParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) WithTimeout(timeout time.Duration) *FindConfigInternodeAuthenticatorParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) WithContext(ctx context.Context) *FindConfigInternodeAuthenticatorParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) WithHTTPClient(client *http.Client) *FindConfigInternodeAuthenticatorParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the find config internode authenticator params
func (o *FindConfigInternodeAuthenticatorParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *FindConfigInternodeAuthenticatorParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
