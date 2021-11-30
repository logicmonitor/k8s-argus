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

// DeleteDeviceByIDReader is a Reader for the DeleteDeviceByID structure.
type DeleteDeviceByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDeviceByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteDeviceByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewDeleteDeviceByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteDeviceByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteDeviceByIDOK creates a DeleteDeviceByIDOK with default headers values
func NewDeleteDeviceByIDOK() *DeleteDeviceByIDOK {
	return &DeleteDeviceByIDOK{}
}

/* DeleteDeviceByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type DeleteDeviceByIDOK struct {
	Payload interface{}
}

func (o *DeleteDeviceByIDOK) Error() string {
	return fmt.Sprintf("[DELETE /device/devices/{id}][%d] deleteDeviceByIdOK  %+v", 200, o.Payload)
}
func (o *DeleteDeviceByIDOK) GetPayload() interface{} {
	return o.Payload
}

func (o *DeleteDeviceByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteDeviceByIDTooManyRequests creates a DeleteDeviceByIDTooManyRequests with default headers values
func NewDeleteDeviceByIDTooManyRequests() *DeleteDeviceByIDTooManyRequests {
	return &DeleteDeviceByIDTooManyRequests{}
}

/* DeleteDeviceByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type DeleteDeviceByIDTooManyRequests struct {

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

func (o *DeleteDeviceByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[DELETE /device/devices/{id}][%d] deleteDeviceByIdTooManyRequests ", 429)
}

func (o *DeleteDeviceByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewDeleteDeviceByIDDefault creates a DeleteDeviceByIDDefault with default headers values
func NewDeleteDeviceByIDDefault(code int) *DeleteDeviceByIDDefault {
	return &DeleteDeviceByIDDefault{
		_statusCode: code,
	}
}

/* DeleteDeviceByIDDefault describes a response with status code -1, with default header values.

Error
*/
type DeleteDeviceByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete device by Id default response
func (o *DeleteDeviceByIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteDeviceByIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /device/devices/{id}][%d] deleteDeviceById default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteDeviceByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteDeviceByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
