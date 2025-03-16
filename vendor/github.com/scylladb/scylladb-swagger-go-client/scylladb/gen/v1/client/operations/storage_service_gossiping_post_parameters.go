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

// NewStorageServiceGossipingPostParams creates a new StorageServiceGossipingPostParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewStorageServiceGossipingPostParams() *StorageServiceGossipingPostParams {
	return &StorageServiceGossipingPostParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewStorageServiceGossipingPostParamsWithTimeout creates a new StorageServiceGossipingPostParams object
// with the ability to set a timeout on a request.
func NewStorageServiceGossipingPostParamsWithTimeout(timeout time.Duration) *StorageServiceGossipingPostParams {
	return &StorageServiceGossipingPostParams{
		timeout: timeout,
	}
}

// NewStorageServiceGossipingPostParamsWithContext creates a new StorageServiceGossipingPostParams object
// with the ability to set a context for a request.
func NewStorageServiceGossipingPostParamsWithContext(ctx context.Context) *StorageServiceGossipingPostParams {
	return &StorageServiceGossipingPostParams{
		Context: ctx,
	}
}

// NewStorageServiceGossipingPostParamsWithHTTPClient creates a new StorageServiceGossipingPostParams object
// with the ability to set a custom HTTPClient for a request.
func NewStorageServiceGossipingPostParamsWithHTTPClient(client *http.Client) *StorageServiceGossipingPostParams {
	return &StorageServiceGossipingPostParams{
		HTTPClient: client,
	}
}

/*
StorageServiceGossipingPostParams contains all the parameters to send to the API endpoint

	for the storage service gossiping post operation.

	Typically these are written to a http.Request.
*/
type StorageServiceGossipingPostParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the storage service gossiping post params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceGossipingPostParams) WithDefaults() *StorageServiceGossipingPostParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the storage service gossiping post params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceGossipingPostParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) WithTimeout(timeout time.Duration) *StorageServiceGossipingPostParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) WithContext(ctx context.Context) *StorageServiceGossipingPostParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) WithHTTPClient(client *http.Client) *StorageServiceGossipingPostParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the storage service gossiping post params
func (o *StorageServiceGossipingPostParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *StorageServiceGossipingPostParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
