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

// AddAlertNoteByIDReader is a Reader for the AddAlertNoteByID structure.
type AddAlertNoteByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddAlertNoteByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewAddAlertNoteByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewAddAlertNoteByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddAlertNoteByIDOK creates a AddAlertNoteByIDOK with default headers values
func NewAddAlertNoteByIDOK() *AddAlertNoteByIDOK {
	return &AddAlertNoteByIDOK{}
}

/*AddAlertNoteByIDOK handles this case with default header values.

successful operation
*/
type AddAlertNoteByIDOK struct {
	Payload interface{}
}

func (o *AddAlertNoteByIDOK) Error() string {
	return fmt.Sprintf("[POST /alert/alerts/{id}/note][%d] addAlertNoteByIdOK  %+v", 200, o.Payload)
}

func (o *AddAlertNoteByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddAlertNoteByIDDefault creates a AddAlertNoteByIDDefault with default headers values
func NewAddAlertNoteByIDDefault(code int) *AddAlertNoteByIDDefault {
	return &AddAlertNoteByIDDefault{
		_statusCode: code,
	}
}

/*AddAlertNoteByIDDefault handles this case with default header values.

Error
*/
type AddAlertNoteByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the add alert note by Id default response
func (o *AddAlertNoteByIDDefault) Code() int {
	return o._statusCode
}

func (o *AddAlertNoteByIDDefault) Error() string {
	return fmt.Sprintf("[POST /alert/alerts/{id}/note][%d] addAlertNoteById default  %+v", o._statusCode, o.Payload)
}

func (o *AddAlertNoteByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
