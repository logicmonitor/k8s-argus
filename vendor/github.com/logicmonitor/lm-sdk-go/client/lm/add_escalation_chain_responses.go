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

// AddEscalationChainReader is a Reader for the AddEscalationChain structure.
type AddEscalationChainReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddEscalationChainReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAddEscalationChainOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewAddEscalationChainTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewAddEscalationChainDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddEscalationChainOK creates a AddEscalationChainOK with default headers values
func NewAddEscalationChainOK() *AddEscalationChainOK {
	return &AddEscalationChainOK{}
}

/* AddEscalationChainOK describes a response with status code 200, with default header values.

successful operation
*/
type AddEscalationChainOK struct {
	Payload *models.EscalatingChain
}

func (o *AddEscalationChainOK) Error() string {
	return fmt.Sprintf("[POST /setting/alert/chains][%d] addEscalationChainOK  %+v", 200, o.Payload)
}
func (o *AddEscalationChainOK) GetPayload() *models.EscalatingChain {
	return o.Payload
}

func (o *AddEscalationChainOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EscalatingChain)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddEscalationChainTooManyRequests creates a AddEscalationChainTooManyRequests with default headers values
func NewAddEscalationChainTooManyRequests() *AddEscalationChainTooManyRequests {
	return &AddEscalationChainTooManyRequests{}
}

/* AddEscalationChainTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type AddEscalationChainTooManyRequests struct {

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

func (o *AddEscalationChainTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /setting/alert/chains][%d] addEscalationChainTooManyRequests ", 429)
}

func (o *AddEscalationChainTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewAddEscalationChainDefault creates a AddEscalationChainDefault with default headers values
func NewAddEscalationChainDefault(code int) *AddEscalationChainDefault {
	return &AddEscalationChainDefault{
		_statusCode: code,
	}
}

/* AddEscalationChainDefault describes a response with status code -1, with default header values.

Error
*/
type AddEscalationChainDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the add escalation chain default response
func (o *AddEscalationChainDefault) Code() int {
	return o._statusCode
}

func (o *AddEscalationChainDefault) Error() string {
	return fmt.Sprintf("[POST /setting/alert/chains][%d] addEscalationChain default  %+v", o._statusCode, o.Payload)
}
func (o *AddEscalationChainDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddEscalationChainDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
