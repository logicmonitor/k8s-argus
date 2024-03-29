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

// PatchDashboardGroupByIDReader is a Reader for the PatchDashboardGroupByID structure.
type PatchDashboardGroupByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchDashboardGroupByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchDashboardGroupByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewPatchDashboardGroupByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPatchDashboardGroupByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPatchDashboardGroupByIDOK creates a PatchDashboardGroupByIDOK with default headers values
func NewPatchDashboardGroupByIDOK() *PatchDashboardGroupByIDOK {
	return &PatchDashboardGroupByIDOK{}
}

/* PatchDashboardGroupByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type PatchDashboardGroupByIDOK struct {
	Payload *models.DashboardGroup
}

func (o *PatchDashboardGroupByIDOK) Error() string {
	return fmt.Sprintf("[PATCH /dashboard/groups/{id}][%d] patchDashboardGroupByIdOK  %+v", 200, o.Payload)
}
func (o *PatchDashboardGroupByIDOK) GetPayload() *models.DashboardGroup {
	return o.Payload
}

func (o *PatchDashboardGroupByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DashboardGroup)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchDashboardGroupByIDTooManyRequests creates a PatchDashboardGroupByIDTooManyRequests with default headers values
func NewPatchDashboardGroupByIDTooManyRequests() *PatchDashboardGroupByIDTooManyRequests {
	return &PatchDashboardGroupByIDTooManyRequests{}
}

/* PatchDashboardGroupByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type PatchDashboardGroupByIDTooManyRequests struct {

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

func (o *PatchDashboardGroupByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[PATCH /dashboard/groups/{id}][%d] patchDashboardGroupByIdTooManyRequests ", 429)
}

func (o *PatchDashboardGroupByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewPatchDashboardGroupByIDDefault creates a PatchDashboardGroupByIDDefault with default headers values
func NewPatchDashboardGroupByIDDefault(code int) *PatchDashboardGroupByIDDefault {
	return &PatchDashboardGroupByIDDefault{
		_statusCode: code,
	}
}

/* PatchDashboardGroupByIDDefault describes a response with status code -1, with default header values.

Error
*/
type PatchDashboardGroupByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the patch dashboard group by Id default response
func (o *PatchDashboardGroupByIDDefault) Code() int {
	return o._statusCode
}

func (o *PatchDashboardGroupByIDDefault) Error() string {
	return fmt.Sprintf("[PATCH /dashboard/groups/{id}][%d] patchDashboardGroupById default  %+v", o._statusCode, o.Payload)
}
func (o *PatchDashboardGroupByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PatchDashboardGroupByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
