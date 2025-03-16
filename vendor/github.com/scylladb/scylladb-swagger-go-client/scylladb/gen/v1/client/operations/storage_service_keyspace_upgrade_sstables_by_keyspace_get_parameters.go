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

// NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams creates a new StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams() *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	return &StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithTimeout creates a new StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams object
// with the ability to set a timeout on a request.
func NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithTimeout(timeout time.Duration) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	return &StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams{
		timeout: timeout,
	}
}

// NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithContext creates a new StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams object
// with the ability to set a context for a request.
func NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithContext(ctx context.Context) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	return &StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams{
		Context: ctx,
	}
}

// NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithHTTPClient creates a new StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewStorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParamsWithHTTPClient(client *http.Client) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	return &StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams{
		HTTPClient: client,
	}
}

/*
StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams contains all the parameters to send to the API endpoint

	for the storage service keyspace upgrade sstables by keyspace get operation.

	Typically these are written to a http.Request.
*/
type StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams struct {

	/* Cf.

	   Comma seperated column family names
	*/
	Cf *string

	/* ExcludeCurrentVersion.

	   When set to true exclude current version
	*/
	ExcludeCurrentVersion *bool

	/* Keyspace.

	   The keyspace
	*/
	Keyspace string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the storage service keyspace upgrade sstables by keyspace get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithDefaults() *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the storage service keyspace upgrade sstables by keyspace get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithTimeout(timeout time.Duration) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithContext(ctx context.Context) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithHTTPClient(client *http.Client) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCf adds the cf to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithCf(cf *string) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetCf(cf)
	return o
}

// SetCf adds the cf to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetCf(cf *string) {
	o.Cf = cf
}

// WithExcludeCurrentVersion adds the excludeCurrentVersion to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithExcludeCurrentVersion(excludeCurrentVersion *bool) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetExcludeCurrentVersion(excludeCurrentVersion)
	return o
}

// SetExcludeCurrentVersion adds the excludeCurrentVersion to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetExcludeCurrentVersion(excludeCurrentVersion *bool) {
	o.ExcludeCurrentVersion = excludeCurrentVersion
}

// WithKeyspace adds the keyspace to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WithKeyspace(keyspace string) *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams {
	o.SetKeyspace(keyspace)
	return o
}

// SetKeyspace adds the keyspace to the storage service keyspace upgrade sstables by keyspace get params
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) SetKeyspace(keyspace string) {
	o.Keyspace = keyspace
}

// WriteToRequest writes these params to a swagger request
func (o *StorageServiceKeyspaceUpgradeSstablesByKeyspaceGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Cf != nil {

		// query param cf
		var qrCf string

		if o.Cf != nil {
			qrCf = *o.Cf
		}
		qCf := qrCf
		if qCf != "" {

			if err := r.SetQueryParam("cf", qCf); err != nil {
				return err
			}
		}
	}

	if o.ExcludeCurrentVersion != nil {

		// query param exclude_current_version
		var qrExcludeCurrentVersion bool

		if o.ExcludeCurrentVersion != nil {
			qrExcludeCurrentVersion = *o.ExcludeCurrentVersion
		}
		qExcludeCurrentVersion := swag.FormatBool(qrExcludeCurrentVersion)
		if qExcludeCurrentVersion != "" {

			if err := r.SetQueryParam("exclude_current_version", qExcludeCurrentVersion); err != nil {
				return err
			}
		}
	}

	// path param keyspace
	if err := r.SetPathParam("keyspace", o.Keyspace); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
