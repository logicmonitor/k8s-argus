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

// DeleteAppliesToFunctionByIDReader is a Reader for the DeleteAppliesToFunctionByID structure.
type DeleteAppliesToFunctionByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteAppliesToFunctionByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteAppliesToFunctionByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewDeleteAppliesToFunctionByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteAppliesToFunctionByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteAppliesToFunctionByIDOK creates a DeleteAppliesToFunctionByIDOK with default headers values
func NewDeleteAppliesToFunctionByIDOK() *DeleteAppliesToFunctionByIDOK {
	return &DeleteAppliesToFunctionByIDOK{}
}

/* DeleteAppliesToFunctionByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type DeleteAppliesToFunctionByIDOK struct {
	Payload interface{}
}

func (o *DeleteAppliesToFunctionByIDOK) Error() string {
	return fmt.Sprintf("[DELETE /setting/functions/{id}][%d] deleteAppliesToFunctionByIdOK  %+v", 200, o.Payload)
}
func (o *DeleteAppliesToFunctionByIDOK) GetPayload() interface{} {
	return o.Payload
}

func (o *DeleteAppliesToFunctionByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteAppliesToFunctionByIDTooManyRequests creates a DeleteAppliesToFunctionByIDTooManyRequests with default headers values
func NewDeleteAppliesToFunctionByIDTooManyRequests() *DeleteAppliesToFunctionByIDTooManyRequests {
	return &DeleteAppliesToFunctionByIDTooManyRequests{}
}

/* DeleteAppliesToFunctionByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type DeleteAppliesToFunctionByIDTooManyRequests struct {

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

func (o *DeleteAppliesToFunctionByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[DELETE /setting/functions/{id}][%d] deleteAppliesToFunctionByIdTooManyRequests ", 429)
}

func (o *DeleteAppliesToFunctionByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewDeleteAppliesToFunctionByIDDefault creates a DeleteAppliesToFunctionByIDDefault with default headers values
func NewDeleteAppliesToFunctionByIDDefault(code int) *DeleteAppliesToFunctionByIDDefault {
	return &DeleteAppliesToFunctionByIDDefault{
		_statusCode: code,
	}
}

/* DeleteAppliesToFunctionByIDDefault describes a response with status code -1, with default header values.

Error
*/
type DeleteAppliesToFunctionByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete applies to function by Id default response
func (o *DeleteAppliesToFunctionByIDDefault) Code() int {
	return o._statusCode
}

func (o *DeleteAppliesToFunctionByIDDefault) Error() string {
	return fmt.Sprintf("[DELETE /setting/functions/{id}][%d] deleteAppliesToFunctionById default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteAppliesToFunctionByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteAppliesToFunctionByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
