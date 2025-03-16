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

// StorageProxyCounterWriteRPCTimeoutPostReader is a Reader for the StorageProxyCounterWriteRPCTimeoutPost structure.
type StorageProxyCounterWriteRPCTimeoutPostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageProxyCounterWriteRPCTimeoutPostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageProxyCounterWriteRPCTimeoutPostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageProxyCounterWriteRPCTimeoutPostDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageProxyCounterWriteRPCTimeoutPostOK creates a StorageProxyCounterWriteRPCTimeoutPostOK with default headers values
func NewStorageProxyCounterWriteRPCTimeoutPostOK() *StorageProxyCounterWriteRPCTimeoutPostOK {
	return &StorageProxyCounterWriteRPCTimeoutPostOK{}
}

/*
StorageProxyCounterWriteRPCTimeoutPostOK handles this case with default header values.

StorageProxyCounterWriteRPCTimeoutPostOK storage proxy counter write Rpc timeout post o k
*/
type StorageProxyCounterWriteRPCTimeoutPostOK struct {
}

func (o *StorageProxyCounterWriteRPCTimeoutPostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewStorageProxyCounterWriteRPCTimeoutPostDefault creates a StorageProxyCounterWriteRPCTimeoutPostDefault with default headers values
func NewStorageProxyCounterWriteRPCTimeoutPostDefault(code int) *StorageProxyCounterWriteRPCTimeoutPostDefault {
	return &StorageProxyCounterWriteRPCTimeoutPostDefault{
		_statusCode: code,
	}
}

/*
StorageProxyCounterWriteRPCTimeoutPostDefault handles this case with default header values.

internal server error
*/
type StorageProxyCounterWriteRPCTimeoutPostDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage proxy counter write Rpc timeout post default response
func (o *StorageProxyCounterWriteRPCTimeoutPostDefault) Code() int {
	return o._statusCode
}

func (o *StorageProxyCounterWriteRPCTimeoutPostDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageProxyCounterWriteRPCTimeoutPostDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageProxyCounterWriteRPCTimeoutPostDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
