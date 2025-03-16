// Code generated by go-swagger; DO NOT EDIT.

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/scylladb/scylladb-swagger-go-client/scylladb/gen/v2/models"
)

// FindConfigViewHintsDirectoryReader is a Reader for the FindConfigViewHintsDirectory structure.
type FindConfigViewHintsDirectoryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FindConfigViewHintsDirectoryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFindConfigViewHintsDirectoryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFindConfigViewHintsDirectoryDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFindConfigViewHintsDirectoryOK creates a FindConfigViewHintsDirectoryOK with default headers values
func NewFindConfigViewHintsDirectoryOK() *FindConfigViewHintsDirectoryOK {
	return &FindConfigViewHintsDirectoryOK{}
}

/*
FindConfigViewHintsDirectoryOK handles this case with default header values.

Config value
*/
type FindConfigViewHintsDirectoryOK struct {
	Payload string
}

func (o *FindConfigViewHintsDirectoryOK) GetPayload() string {
	return o.Payload
}

func (o *FindConfigViewHintsDirectoryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindConfigViewHintsDirectoryDefault creates a FindConfigViewHintsDirectoryDefault with default headers values
func NewFindConfigViewHintsDirectoryDefault(code int) *FindConfigViewHintsDirectoryDefault {
	return &FindConfigViewHintsDirectoryDefault{
		_statusCode: code,
	}
}

/*
FindConfigViewHintsDirectoryDefault handles this case with default header values.

unexpected error
*/
type FindConfigViewHintsDirectoryDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the find config view hints directory default response
func (o *FindConfigViewHintsDirectoryDefault) Code() int {
	return o._statusCode
}

func (o *FindConfigViewHintsDirectoryDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *FindConfigViewHintsDirectoryDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *FindConfigViewHintsDirectoryDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
