// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/scylladb/scylladb-swagger-go-client/scylladb/gen/v1/models"
)

// StorageProxyMetricsReadTimeoutsGetReader is a Reader for the StorageProxyMetricsReadTimeoutsGet structure.
type StorageProxyMetricsReadTimeoutsGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageProxyMetricsReadTimeoutsGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageProxyMetricsReadTimeoutsGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageProxyMetricsReadTimeoutsGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageProxyMetricsReadTimeoutsGetOK creates a StorageProxyMetricsReadTimeoutsGetOK with default headers values
func NewStorageProxyMetricsReadTimeoutsGetOK() *StorageProxyMetricsReadTimeoutsGetOK {
	return &StorageProxyMetricsReadTimeoutsGetOK{}
}

/*
StorageProxyMetricsReadTimeoutsGetOK handles this case with default header values.

StorageProxyMetricsReadTimeoutsGetOK storage proxy metrics read timeouts get o k
*/
type StorageProxyMetricsReadTimeoutsGetOK struct {
	Payload int32
}

func (o *StorageProxyMetricsReadTimeoutsGetOK) GetPayload() int32 {
	return o.Payload
}

func (o *StorageProxyMetricsReadTimeoutsGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStorageProxyMetricsReadTimeoutsGetDefault creates a StorageProxyMetricsReadTimeoutsGetDefault with default headers values
func NewStorageProxyMetricsReadTimeoutsGetDefault(code int) *StorageProxyMetricsReadTimeoutsGetDefault {
	return &StorageProxyMetricsReadTimeoutsGetDefault{
		_statusCode: code,
	}
}

/*
StorageProxyMetricsReadTimeoutsGetDefault handles this case with default header values.

internal server error
*/
type StorageProxyMetricsReadTimeoutsGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage proxy metrics read timeouts get default response
func (o *StorageProxyMetricsReadTimeoutsGetDefault) Code() int {
	return o._statusCode
}

func (o *StorageProxyMetricsReadTimeoutsGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageProxyMetricsReadTimeoutsGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageProxyMetricsReadTimeoutsGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
