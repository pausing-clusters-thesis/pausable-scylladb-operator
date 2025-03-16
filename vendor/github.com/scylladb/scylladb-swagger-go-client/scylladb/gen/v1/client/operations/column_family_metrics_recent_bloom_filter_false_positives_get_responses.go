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

// ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetReader is a Reader for the ColumnFamilyMetricsRecentBloomFilterFalsePositivesGet structure.
type ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK creates a ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK with default headers values
func NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK() *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK {
	return &ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK{}
}

/*
ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK handles this case with default header values.

ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK column family metrics recent bloom filter false positives get o k
*/
type ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault creates a ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault with default headers values
func NewColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault(code int) *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault {
	return &ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics recent bloom filter false positives get default response
func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsRecentBloomFilterFalsePositivesGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
