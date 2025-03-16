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

// FindConfigAuthenticatorReader is a Reader for the FindConfigAuthenticator structure.
type FindConfigAuthenticatorReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FindConfigAuthenticatorReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewFindConfigAuthenticatorOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewFindConfigAuthenticatorDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFindConfigAuthenticatorOK creates a FindConfigAuthenticatorOK with default headers values
func NewFindConfigAuthenticatorOK() *FindConfigAuthenticatorOK {
	return &FindConfigAuthenticatorOK{}
}

/*
FindConfigAuthenticatorOK handles this case with default header values.

Config value
*/
type FindConfigAuthenticatorOK struct {
	Payload string
}

func (o *FindConfigAuthenticatorOK) GetPayload() string {
	return o.Payload
}

func (o *FindConfigAuthenticatorOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindConfigAuthenticatorDefault creates a FindConfigAuthenticatorDefault with default headers values
func NewFindConfigAuthenticatorDefault(code int) *FindConfigAuthenticatorDefault {
	return &FindConfigAuthenticatorDefault{
		_statusCode: code,
	}
}

/*
FindConfigAuthenticatorDefault handles this case with default header values.

unexpected error
*/
type FindConfigAuthenticatorDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the find config authenticator default response
func (o *FindConfigAuthenticatorDefault) Code() int {
	return o._statusCode
}

func (o *FindConfigAuthenticatorDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *FindConfigAuthenticatorDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *FindConfigAuthenticatorDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
