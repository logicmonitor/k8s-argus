// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/vkumbhar94/lm-sdk-go/models"
)

// DeleteWebsiteGroupByIDReader is a Reader for the DeleteWebsiteGroupByID structure.
type DeleteWebsiteGroupByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteWebsiteGroupByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewDeleteWebsiteGroupByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewDeleteWebsiteGroupByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteWebsiteGroupByIDOK creates a DeleteWebsiteGroupByIDOK with default headers values
func NewDeleteWebsiteGroupByIDOK() *DeleteWebsiteGroupByIDOK {
	return &DeleteWebsiteGroupByIDOK{}
}

/*DeleteWebsiteGroupByIDOK handles this case with default header values.

successful operation
*/
type DeleteWebsiteGroupByIDOK struct {
	Payload interface{}
}

func (o *DeleteWebsiteGroupByIDOK) Error() string {
	return fmt.Sprintf("[DELETE /website/groups/{id}][%d] deleteWebsiteGroupByIdOK  %+v", 200, o.Payload)
}

func (o *DeleteWebsiteGroupByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteWebsiteGroupByIDDefault creates a DeleteWebsiteGroupByIDDefault with default headers values
func NewDeleteWebsiteGroupByIDDefault(code int) *DeleteWebsiteGroupByIDDefault {
	return &DeleteWebsiteGroupByIDDefault{
		_statusCode: code,
	}
}

/*DeleteWebsiteGroupByIDDefault handles this case with default header values.

Error
*/
type DeleteWebsiteGroupByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete website group by Id default response
func (o *DeleteWebsiteGroupByIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteWebsiteGroupByIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /website/groups/{id}][%d] deleteWebsiteGroupById default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteWebsiteGroupByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}