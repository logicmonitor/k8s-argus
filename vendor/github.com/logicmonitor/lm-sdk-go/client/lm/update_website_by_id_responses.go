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

// UpdateWebsiteByIDReader is a Reader for the UpdateWebsiteByID structure.
type UpdateWebsiteByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateWebsiteByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateWebsiteByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewUpdateWebsiteByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateWebsiteByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateWebsiteByIDOK creates a UpdateWebsiteByIDOK with default headers values
func NewUpdateWebsiteByIDOK() *UpdateWebsiteByIDOK {
	return &UpdateWebsiteByIDOK{}
}

/* UpdateWebsiteByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type UpdateWebsiteByIDOK struct {
	Payload models.Website
}

func (o *UpdateWebsiteByIDOK) Error() string {
	return fmt.Sprintf("[PUT /website/websites/{id}][%d] updateWebsiteByIdOK  %+v", 200, o.Payload)
}
func (o *UpdateWebsiteByIDOK) GetPayload() models.Website {
	return o.Payload
}

func (o *UpdateWebsiteByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload as interface type
	payload, err := models.UnmarshalWebsite(response.Body(), consumer)
	if err != nil {
		return err
	}
	o.Payload = payload

	return nil
}

// NewUpdateWebsiteByIDTooManyRequests creates a UpdateWebsiteByIDTooManyRequests with default headers values
func NewUpdateWebsiteByIDTooManyRequests() *UpdateWebsiteByIDTooManyRequests {
	return &UpdateWebsiteByIDTooManyRequests{}
}

/* UpdateWebsiteByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type UpdateWebsiteByIDTooManyRequests struct {

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

func (o *UpdateWebsiteByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[PUT /website/websites/{id}][%d] updateWebsiteByIdTooManyRequests ", 429)
}

func (o *UpdateWebsiteByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewUpdateWebsiteByIDDefault creates a UpdateWebsiteByIDDefault with default headers values
func NewUpdateWebsiteByIDDefault(code int) *UpdateWebsiteByIDDefault {
	return &UpdateWebsiteByIDDefault{
		_statusCode: code,
	}
}

/* UpdateWebsiteByIDDefault describes a response with status code -1, with default header values.

Error
*/
type UpdateWebsiteByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the update website by Id default response
func (o *UpdateWebsiteByIDDefault) Code() int {
	return o._statusCode
}

func (o *UpdateWebsiteByIDDefault) Error() string {
	return fmt.Sprintf("[PUT /website/websites/{id}][%d] updateWebsiteById default  %+v", o._statusCode, o.Payload)
}
func (o *UpdateWebsiteByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateWebsiteByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
