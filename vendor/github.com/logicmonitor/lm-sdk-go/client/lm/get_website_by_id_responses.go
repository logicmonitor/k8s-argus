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

// GetWebsiteByIDReader is a Reader for the GetWebsiteByID structure.
type GetWebsiteByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWebsiteByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWebsiteByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetWebsiteByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetWebsiteByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetWebsiteByIDOK creates a GetWebsiteByIDOK with default headers values
func NewGetWebsiteByIDOK() *GetWebsiteByIDOK {
	return &GetWebsiteByIDOK{}
}

/* GetWebsiteByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type GetWebsiteByIDOK struct {
	Payload models.Website
}

func (o *GetWebsiteByIDOK) Error() string {
	return fmt.Sprintf("[GET /website/websites/{id}][%d] getWebsiteByIdOK  %+v", 200, o.Payload)
}
func (o *GetWebsiteByIDOK) GetPayload() models.Website {
	return o.Payload
}

func (o *GetWebsiteByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload as interface type
	payload, err := models.UnmarshalWebsite(response.Body(), consumer)
	if err != nil {
		return err
	}
	o.Payload = payload

	return nil
}

// NewGetWebsiteByIDTooManyRequests creates a GetWebsiteByIDTooManyRequests with default headers values
func NewGetWebsiteByIDTooManyRequests() *GetWebsiteByIDTooManyRequests {
	return &GetWebsiteByIDTooManyRequests{}
}

/* GetWebsiteByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetWebsiteByIDTooManyRequests struct {

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

func (o *GetWebsiteByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /website/websites/{id}][%d] getWebsiteByIdTooManyRequests ", 429)
}

func (o *GetWebsiteByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetWebsiteByIDDefault creates a GetWebsiteByIDDefault with default headers values
func NewGetWebsiteByIDDefault(code int) *GetWebsiteByIDDefault {
	return &GetWebsiteByIDDefault{
		_statusCode: code,
	}
}

/* GetWebsiteByIDDefault describes a response with status code -1, with default header values.

Error
*/
type GetWebsiteByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get website by Id default response
func (o *GetWebsiteByIDDefault) Code() int {
	return o._statusCode
}

func (o *GetWebsiteByIDDefault) Error() string {
	return fmt.Sprintf("[GET /website/websites/{id}][%d] getWebsiteById default  %+v", o._statusCode, o.Payload)
}
func (o *GetWebsiteByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetWebsiteByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
