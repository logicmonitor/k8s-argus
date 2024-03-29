// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/logicmonitor/lm-sdk-go/models"
)

// DeleteDeviceGroupByIDReader is a Reader for the DeleteDeviceGroupByID structure.
type DeleteDeviceGroupByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDeviceGroupByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteDeviceGroupByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewDeleteDeviceGroupByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteDeviceGroupByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteDeviceGroupByIDOK creates a DeleteDeviceGroupByIDOK with default headers values
func NewDeleteDeviceGroupByIDOK() *DeleteDeviceGroupByIDOK {
	return &DeleteDeviceGroupByIDOK{}
}

/* DeleteDeviceGroupByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type DeleteDeviceGroupByIDOK struct {
	Payload interface{}
}

func (o *DeleteDeviceGroupByIDOK) Error() string {
	return fmt.Sprintf("[DELETE /device/groups/{id}][%d] deleteDeviceGroupByIdOK  %+v", 200, o.Payload)
}
func (o *DeleteDeviceGroupByIDOK) GetPayload() interface{} {
	return o.Payload
}

func (o *DeleteDeviceGroupByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteDeviceGroupByIDTooManyRequests creates a DeleteDeviceGroupByIDTooManyRequests with default headers values
func NewDeleteDeviceGroupByIDTooManyRequests() *DeleteDeviceGroupByIDTooManyRequests {
	return &DeleteDeviceGroupByIDTooManyRequests{}
}

/* DeleteDeviceGroupByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type DeleteDeviceGroupByIDTooManyRequests struct {

	/* Request limit per X-Rate-Limit-Window
	 */
	XRateLimitLimit int64

	/* The number of requests left for the time window
	 */
	XRateLimitRemaining int64

	/* The rolling time window length with the unit of second
	 */
	XRateLimitWindow int64
}

func (o *DeleteDeviceGroupByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[DELETE /device/groups/{id}][%d] deleteDeviceGroupByIdTooManyRequests ", 429)
}

func (o *DeleteDeviceGroupByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header x-rate-limit-limit
	hdrXRateLimitLimit := response.GetHeader("x-rate-limit-limit")

	if hdrXRateLimitLimit != "" {
		valxRateLimitLimit, err := swag.ConvertInt64(hdrXRateLimitLimit)
		if err != nil {
			return errors.InvalidType("x-rate-limit-limit", "header", "int64", hdrXRateLimitLimit)
		}
		o.XRateLimitLimit = valxRateLimitLimit
	}

	// hydrates response header x-rate-limit-remaining
	hdrXRateLimitRemaining := response.GetHeader("x-rate-limit-remaining")

	if hdrXRateLimitRemaining != "" {
		valxRateLimitRemaining, err := swag.ConvertInt64(hdrXRateLimitRemaining)
		if err != nil {
			return errors.InvalidType("x-rate-limit-remaining", "header", "int64", hdrXRateLimitRemaining)
		}
		o.XRateLimitRemaining = valxRateLimitRemaining
	}

	// hydrates response header x-rate-limit-window
	hdrXRateLimitWindow := response.GetHeader("x-rate-limit-window")

	if hdrXRateLimitWindow != "" {
		valxRateLimitWindow, err := swag.ConvertInt64(hdrXRateLimitWindow)
		if err != nil {
			return errors.InvalidType("x-rate-limit-window", "header", "int64", hdrXRateLimitWindow)
		}
		o.XRateLimitWindow = valxRateLimitWindow
	}

	return nil
}

// NewDeleteDeviceGroupByIDDefault creates a DeleteDeviceGroupByIDDefault with default headers values
func NewDeleteDeviceGroupByIDDefault(code int) *DeleteDeviceGroupByIDDefault {
	return &DeleteDeviceGroupByIDDefault{
		_statusCode: code,
	}
}

/* DeleteDeviceGroupByIDDefault describes a response with status code -1, with default header values.

Error
*/
type DeleteDeviceGroupByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete device group by Id default response
func (o *DeleteDeviceGroupByIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteDeviceGroupByIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /device/groups/{id}][%d] deleteDeviceGroupById default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteDeviceGroupByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteDeviceGroupByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
