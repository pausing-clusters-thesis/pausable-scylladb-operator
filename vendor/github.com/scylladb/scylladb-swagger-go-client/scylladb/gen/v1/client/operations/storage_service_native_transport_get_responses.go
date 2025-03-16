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

// StorageServiceNativeTransportGetReader is a Reader for the StorageServiceNativeTransportGet structure.
type StorageServiceNativeTransportGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageServiceNativeTransportGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageServiceNativeTransportGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageServiceNativeTransportGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageServiceNativeTransportGetOK creates a StorageServiceNativeTransportGetOK with default headers values
func NewStorageServiceNativeTransportGetOK() *StorageServiceNativeTransportGetOK {
	return &StorageServiceNativeTransportGetOK{}
}

/*
StorageServiceNativeTransportGetOK handles this case with default header values.

StorageServiceNativeTransportGetOK storage service native transport get o k
*/
type StorageServiceNativeTransportGetOK struct {
	Payload bool
}

func (o *StorageServiceNativeTransportGetOK) GetPayload() bool {
	return o.Payload
}

func (o *StorageServiceNativeTransportGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStorageServiceNativeTransportGetDefault creates a StorageServiceNativeTransportGetDefault with default headers values
func NewStorageServiceNativeTransportGetDefault(code int) *StorageServiceNativeTransportGetDefault {
	return &StorageServiceNativeTransportGetDefault{
		_statusCode: code,
	}
}

/*
StorageServiceNativeTransportGetDefault handles this case with default header values.

internal server error
*/
type StorageServiceNativeTransportGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage service native transport get default response
func (o *StorageServiceNativeTransportGetDefault) Code() int {
	return o._statusCode
}

func (o *StorageServiceNativeTransportGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageServiceNativeTransportGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageServiceNativeTransportGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
