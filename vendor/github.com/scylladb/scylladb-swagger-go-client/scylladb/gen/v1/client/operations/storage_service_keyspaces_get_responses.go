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

// StorageServiceKeyspacesGetReader is a Reader for the StorageServiceKeyspacesGet structure.
type StorageServiceKeyspacesGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageServiceKeyspacesGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageServiceKeyspacesGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageServiceKeyspacesGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageServiceKeyspacesGetOK creates a StorageServiceKeyspacesGetOK with default headers values
func NewStorageServiceKeyspacesGetOK() *StorageServiceKeyspacesGetOK {
	return &StorageServiceKeyspacesGetOK{}
}

/*
StorageServiceKeyspacesGetOK handles this case with default header values.

StorageServiceKeyspacesGetOK storage service keyspaces get o k
*/
type StorageServiceKeyspacesGetOK struct {
	Payload []string
}

func (o *StorageServiceKeyspacesGetOK) GetPayload() []string {
	return o.Payload
}

func (o *StorageServiceKeyspacesGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStorageServiceKeyspacesGetDefault creates a StorageServiceKeyspacesGetDefault with default headers values
func NewStorageServiceKeyspacesGetDefault(code int) *StorageServiceKeyspacesGetDefault {
	return &StorageServiceKeyspacesGetDefault{
		_statusCode: code,
	}
}

/*
StorageServiceKeyspacesGetDefault handles this case with default header values.

internal server error
*/
type StorageServiceKeyspacesGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage service keyspaces get default response
func (o *StorageServiceKeyspacesGetDefault) Code() int {
	return o._statusCode
}

func (o *StorageServiceKeyspacesGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageServiceKeyspacesGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageServiceKeyspacesGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
