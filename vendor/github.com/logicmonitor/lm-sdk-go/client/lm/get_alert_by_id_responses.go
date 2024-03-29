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

// GetAlertByIDReader is a Reader for the GetAlertByID structure.
type GetAlertByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAlertByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAlertByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetAlertByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetAlertByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetAlertByIDOK creates a GetAlertByIDOK with default headers values
func NewGetAlertByIDOK() *GetAlertByIDOK {
	return &GetAlertByIDOK{}
}

/* GetAlertByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type GetAlertByIDOK struct {
	Payload *models.Alert
}

func (o *GetAlertByIDOK) Error() string {
	return fmt.Sprintf("[GET /alert/alerts/{id}][%d] getAlertByIdOK  %+v", 200, o.Payload)
}
func (o *GetAlertByIDOK) GetPayload() *models.Alert {
	return o.Payload
}

func (o *GetAlertByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Alert)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlertByIDTooManyRequests creates a GetAlertByIDTooManyRequests with default headers values
func NewGetAlertByIDTooManyRequests() *GetAlertByIDTooManyRequests {
	return &GetAlertByIDTooManyRequests{}
}

/* GetAlertByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetAlertByIDTooManyRequests struct {

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

func (o *GetAlertByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /alert/alerts/{id}][%d] getAlertByIdTooManyRequests ", 429)
}

func (o *GetAlertByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetAlertByIDDefault creates a GetAlertByIDDefault with default headers values
func NewGetAlertByIDDefault(code int) *GetAlertByIDDefault {
	return &GetAlertByIDDefault{
		_statusCode: code,
	}
}

/* GetAlertByIDDefault describes a response with status code -1, with default header values.

Error
*/
type GetAlertByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get alert by Id default response
func (o *GetAlertByIDDefault) Code() int {
	return o._statusCode
}

func (o *GetAlertByIDDefault) Error() string {
	return fmt.Sprintf("[GET /alert/alerts/{id}][%d] getAlertById default  %+v", o._statusCode, o.Payload)
}
func (o *GetAlertByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetAlertByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
