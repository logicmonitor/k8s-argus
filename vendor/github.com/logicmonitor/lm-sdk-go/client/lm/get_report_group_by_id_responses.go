// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
	models "github.com/logicmonitor/lm-sdk-go/models"
)

// GetReportGroupByIDReader is a Reader for the GetReportGroupByID structure.
type GetReportGroupByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetReportGroupByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetReportGroupByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetReportGroupByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetReportGroupByIDOK creates a GetReportGroupByIDOK with default headers values
func NewGetReportGroupByIDOK() *GetReportGroupByIDOK {
	return &GetReportGroupByIDOK{}
}

/*GetReportGroupByIDOK handles this case with default header values.

successful operation
*/
type GetReportGroupByIDOK struct {
	Payload *models.ReportGroup
}

func (o *GetReportGroupByIDOK) Error() string {
	return fmt.Sprintf("[GET /report/groups/{id}][%d] getReportGroupByIdOK  %+v", 200, o.Payload)
}

func (o *GetReportGroupByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.ReportGroup)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetReportGroupByIDDefault creates a GetReportGroupByIDDefault with default headers values
func NewGetReportGroupByIDDefault(code int) *GetReportGroupByIDDefault {
	return &GetReportGroupByIDDefault{
		_statusCode: code,
	}
}

/*GetReportGroupByIDDefault handles this case with default header values.

Error
*/
type GetReportGroupByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get report group by Id default response
func (o *GetReportGroupByIDDefault) Code() int {
	return o._statusCode
}

func (o *GetReportGroupByIDDefault) Error() string {
	return fmt.Sprintf("[GET /report/groups/{id}][%d] getReportGroupById default  %+v", o._statusCode, o.Payload)
}

func (o *GetReportGroupByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
