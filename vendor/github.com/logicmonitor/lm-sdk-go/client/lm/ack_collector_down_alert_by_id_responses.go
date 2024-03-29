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

// AckCollectorDownAlertByIDReader is a Reader for the AckCollectorDownAlertByID structure.
type AckCollectorDownAlertByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AckCollectorDownAlertByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAckCollectorDownAlertByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewAckCollectorDownAlertByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewAckCollectorDownAlertByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAckCollectorDownAlertByIDOK creates a AckCollectorDownAlertByIDOK with default headers values
func NewAckCollectorDownAlertByIDOK() *AckCollectorDownAlertByIDOK {
	return &AckCollectorDownAlertByIDOK{}
}

/* AckCollectorDownAlertByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type AckCollectorDownAlertByIDOK struct {
	Payload interface{}
}

func (o *AckCollectorDownAlertByIDOK) Error() string {
	return fmt.Sprintf("[POST /setting/collector/collectors/{id}/ackdown][%d] ackCollectorDownAlertByIdOK  %+v", 200, o.Payload)
}
func (o *AckCollectorDownAlertByIDOK) GetPayload() interface{} {
	return o.Payload
}

func (o *AckCollectorDownAlertByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAckCollectorDownAlertByIDTooManyRequests creates a AckCollectorDownAlertByIDTooManyRequests with default headers values
func NewAckCollectorDownAlertByIDTooManyRequests() *AckCollectorDownAlertByIDTooManyRequests {
	return &AckCollectorDownAlertByIDTooManyRequests{}
}

/* AckCollectorDownAlertByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type AckCollectorDownAlertByIDTooManyRequests struct {

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

func (o *AckCollectorDownAlertByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /setting/collector/collectors/{id}/ackdown][%d] ackCollectorDownAlertByIdTooManyRequests ", 429)
}

func (o *AckCollectorDownAlertByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewAckCollectorDownAlertByIDDefault creates a AckCollectorDownAlertByIDDefault with default headers values
func NewAckCollectorDownAlertByIDDefault(code int) *AckCollectorDownAlertByIDDefault {
	return &AckCollectorDownAlertByIDDefault{
		_statusCode: code,
	}
}

/* AckCollectorDownAlertByIDDefault describes a response with status code -1, with default header values.

Error
*/
type AckCollectorDownAlertByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the ack collector down alert by Id default response
func (o *AckCollectorDownAlertByIDDefault) Code() int {
	return o._statusCode
}

func (o *AckCollectorDownAlertByIDDefault) Error() string {
	return fmt.Sprintf("[POST /setting/collector/collectors/{id}/ackdown][%d] ackCollectorDownAlertById default  %+v", o._statusCode, o.Payload)
}
func (o *AckCollectorDownAlertByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AckCollectorDownAlertByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
