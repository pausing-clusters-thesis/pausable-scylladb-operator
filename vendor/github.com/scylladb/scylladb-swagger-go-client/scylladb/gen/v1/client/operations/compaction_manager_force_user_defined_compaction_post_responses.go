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

// CompactionManagerForceUserDefinedCompactionPostReader is a Reader for the CompactionManagerForceUserDefinedCompactionPost structure.
type CompactionManagerForceUserDefinedCompactionPostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CompactionManagerForceUserDefinedCompactionPostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCompactionManagerForceUserDefinedCompactionPostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCompactionManagerForceUserDefinedCompactionPostDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCompactionManagerForceUserDefinedCompactionPostOK creates a CompactionManagerForceUserDefinedCompactionPostOK with default headers values
func NewCompactionManagerForceUserDefinedCompactionPostOK() *CompactionManagerForceUserDefinedCompactionPostOK {
	return &CompactionManagerForceUserDefinedCompactionPostOK{}
}

/*
CompactionManagerForceUserDefinedCompactionPostOK handles this case with default header values.

CompactionManagerForceUserDefinedCompactionPostOK compaction manager force user defined compaction post o k
*/
type CompactionManagerForceUserDefinedCompactionPostOK struct {
}

func (o *CompactionManagerForceUserDefinedCompactionPostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCompactionManagerForceUserDefinedCompactionPostDefault creates a CompactionManagerForceUserDefinedCompactionPostDefault with default headers values
func NewCompactionManagerForceUserDefinedCompactionPostDefault(code int) *CompactionManagerForceUserDefinedCompactionPostDefault {
	return &CompactionManagerForceUserDefinedCompactionPostDefault{
		_statusCode: code,
	}
}

/*
CompactionManagerForceUserDefinedCompactionPostDefault handles this case with default header values.

internal server error
*/
type CompactionManagerForceUserDefinedCompactionPostDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the compaction manager force user defined compaction post default response
func (o *CompactionManagerForceUserDefinedCompactionPostDefault) Code() int {
	return o._statusCode
}

func (o *CompactionManagerForceUserDefinedCompactionPostDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *CompactionManagerForceUserDefinedCompactionPostDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *CompactionManagerForceUserDefinedCompactionPostDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
