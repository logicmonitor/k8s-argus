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

// GetRecipientGroupListReader is a Reader for the GetRecipientGroupList structure.
type GetRecipientGroupListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetRecipientGroupListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetRecipientGroupListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetRecipientGroupListTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetRecipientGroupListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetRecipientGroupListOK creates a GetRecipientGroupListOK with default headers values
func NewGetRecipientGroupListOK() *GetRecipientGroupListOK {
	return &GetRecipientGroupListOK{}
}

/* GetRecipientGroupListOK describes a response with status code 200, with default header values.

successful operation
*/
type GetRecipientGroupListOK struct {
	Payload *models.RecipientGroupPaginationResponse
}

func (o *GetRecipientGroupListOK) Error() string {
	return fmt.Sprintf("[GET /setting/recipientgroups][%d] getRecipientGroupListOK  %+v", 200, o.Payload)
}
func (o *GetRecipientGroupListOK) GetPayload() *models.RecipientGroupPaginationResponse {
	return o.Payload
}

func (o *GetRecipientGroupListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RecipientGroupPaginationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetRecipientGroupListTooManyRequests creates a GetRecipientGroupListTooManyRequests with default headers values
func NewGetRecipientGroupListTooManyRequests() *GetRecipientGroupListTooManyRequests {
	return &GetRecipientGroupListTooManyRequests{}
}

/* GetRecipientGroupListTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetRecipientGroupListTooManyRequests struct {

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

func (o *GetRecipientGroupListTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /setting/recipientgroups][%d] getRecipientGroupListTooManyRequests ", 429)
}

func (o *GetRecipientGroupListTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetRecipientGroupListDefault creates a GetRecipientGroupListDefault with default headers values
func NewGetRecipientGroupListDefault(code int) *GetRecipientGroupListDefault {
	return &GetRecipientGroupListDefault{
		_statusCode: code,
	}
}

/* GetRecipientGroupListDefault describes a response with status code -1, with default header values.

Error
*/
type GetRecipientGroupListDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get recipient group list default response
func (o *GetRecipientGroupListDefault) Code() int {
	return o._statusCode
}

func (o *GetRecipientGroupListDefault) Error() string {
	return fmt.Sprintf("[GET /setting/recipientgroups][%d] getRecipientGroupList default  %+v", o._statusCode, o.Payload)
}
func (o *GetRecipientGroupListDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetRecipientGroupListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
