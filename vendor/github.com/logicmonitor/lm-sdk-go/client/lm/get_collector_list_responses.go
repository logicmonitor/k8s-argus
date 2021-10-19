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

// GetCollectorListReader is a Reader for the GetCollectorList structure.
type GetCollectorListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCollectorListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCollectorListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetCollectorListTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetCollectorListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCollectorListOK creates a GetCollectorListOK with default headers values
func NewGetCollectorListOK() *GetCollectorListOK {
	return &GetCollectorListOK{}
}

/* GetCollectorListOK describes a response with status code 200, with default header values.

successful operation
*/
type GetCollectorListOK struct {
	Payload *models.CollectorPaginationResponse
}

func (o *GetCollectorListOK) Error() string {
	return fmt.Sprintf("[GET /setting/collector/collectors][%d] getCollectorListOK  %+v", 200, o.Payload)
}
func (o *GetCollectorListOK) GetPayload() *models.CollectorPaginationResponse {
	return o.Payload
}

func (o *GetCollectorListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.CollectorPaginationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCollectorListTooManyRequests creates a GetCollectorListTooManyRequests with default headers values
func NewGetCollectorListTooManyRequests() *GetCollectorListTooManyRequests {
	return &GetCollectorListTooManyRequests{}
}

/* GetCollectorListTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetCollectorListTooManyRequests struct {

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

func (o *GetCollectorListTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /setting/collector/collectors][%d] getCollectorListTooManyRequests ", 429)
}

func (o *GetCollectorListTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetCollectorListDefault creates a GetCollectorListDefault with default headers values
func NewGetCollectorListDefault(code int) *GetCollectorListDefault {
	return &GetCollectorListDefault{
		_statusCode: code,
	}
}

/* GetCollectorListDefault describes a response with status code -1, with default header values.

Error
*/
type GetCollectorListDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get collector list default response
func (o *GetCollectorListDefault) Code() int {
	return o._statusCode
}

func (o *GetCollectorListDefault) Error() string {
	return fmt.Sprintf("[GET /setting/collector/collectors][%d] getCollectorList default  %+v", o._statusCode, o.Payload)
}
func (o *GetCollectorListDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetCollectorListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
