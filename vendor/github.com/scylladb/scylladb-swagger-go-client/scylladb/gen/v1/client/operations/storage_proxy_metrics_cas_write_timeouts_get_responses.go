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

// StorageProxyMetricsCasWriteTimeoutsGetReader is a Reader for the StorageProxyMetricsCasWriteTimeoutsGet structure.
type StorageProxyMetricsCasWriteTimeoutsGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageProxyMetricsCasWriteTimeoutsGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageProxyMetricsCasWriteTimeoutsGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageProxyMetricsCasWriteTimeoutsGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageProxyMetricsCasWriteTimeoutsGetOK creates a StorageProxyMetricsCasWriteTimeoutsGetOK with default headers values
func NewStorageProxyMetricsCasWriteTimeoutsGetOK() *StorageProxyMetricsCasWriteTimeoutsGetOK {
	return &StorageProxyMetricsCasWriteTimeoutsGetOK{}
}

/*
StorageProxyMetricsCasWriteTimeoutsGetOK handles this case with default header values.

StorageProxyMetricsCasWriteTimeoutsGetOK storage proxy metrics cas write timeouts get o k
*/
type StorageProxyMetricsCasWriteTimeoutsGetOK struct {
	Payload interface{}
}

func (o *StorageProxyMetricsCasWriteTimeoutsGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *StorageProxyMetricsCasWriteTimeoutsGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStorageProxyMetricsCasWriteTimeoutsGetDefault creates a StorageProxyMetricsCasWriteTimeoutsGetDefault with default headers values
func NewStorageProxyMetricsCasWriteTimeoutsGetDefault(code int) *StorageProxyMetricsCasWriteTimeoutsGetDefault {
	return &StorageProxyMetricsCasWriteTimeoutsGetDefault{
		_statusCode: code,
	}
}

/*
StorageProxyMetricsCasWriteTimeoutsGetDefault handles this case with default header values.

internal server error
*/
type StorageProxyMetricsCasWriteTimeoutsGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage proxy metrics cas write timeouts get default response
func (o *StorageProxyMetricsCasWriteTimeoutsGetDefault) Code() int {
	return o._statusCode
}

func (o *StorageProxyMetricsCasWriteTimeoutsGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageProxyMetricsCasWriteTimeoutsGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageProxyMetricsCasWriteTimeoutsGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
