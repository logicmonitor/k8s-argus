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

// UpdateDeviceGroupPropertyByNameReader is a Reader for the UpdateDeviceGroupPropertyByName structure.
type UpdateDeviceGroupPropertyByNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateDeviceGroupPropertyByNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateDeviceGroupPropertyByNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewUpdateDeviceGroupPropertyByNameTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateDeviceGroupPropertyByNameDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateDeviceGroupPropertyByNameOK creates a UpdateDeviceGroupPropertyByNameOK with default headers values
func NewUpdateDeviceGroupPropertyByNameOK() *UpdateDeviceGroupPropertyByNameOK {
	return &UpdateDeviceGroupPropertyByNameOK{}
}

/* UpdateDeviceGroupPropertyByNameOK describes a response with status code 200, with default header values.

successful operation
*/
type UpdateDeviceGroupPropertyByNameOK struct {
	Payload *models.EntityProperty
}

func (o *UpdateDeviceGroupPropertyByNameOK) Error() string {
	return fmt.Sprintf("[PUT /device/groups/{gid}/properties/{name}][%d] updateDeviceGroupPropertyByNameOK  %+v", 200, o.Payload)
}
func (o *UpdateDeviceGroupPropertyByNameOK) GetPayload() *models.EntityProperty {
	return o.Payload
}

func (o *UpdateDeviceGroupPropertyByNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EntityProperty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateDeviceGroupPropertyByNameTooManyRequests creates a UpdateDeviceGroupPropertyByNameTooManyRequests with default headers values
func NewUpdateDeviceGroupPropertyByNameTooManyRequests() *UpdateDeviceGroupPropertyByNameTooManyRequests {
	return &UpdateDeviceGroupPropertyByNameTooManyRequests{}
}

/* UpdateDeviceGroupPropertyByNameTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type UpdateDeviceGroupPropertyByNameTooManyRequests struct {

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

func (o *UpdateDeviceGroupPropertyByNameTooManyRequests) Error() string {
	return fmt.Sprintf("[PUT /device/groups/{gid}/properties/{name}][%d] updateDeviceGroupPropertyByNameTooManyRequests ", 429)
}

func (o *UpdateDeviceGroupPropertyByNameTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewUpdateDeviceGroupPropertyByNameDefault creates a UpdateDeviceGroupPropertyByNameDefault with default headers values
func NewUpdateDeviceGroupPropertyByNameDefault(code int) *UpdateDeviceGroupPropertyByNameDefault {
	return &UpdateDeviceGroupPropertyByNameDefault{
		_statusCode: code,
	}
}

/* UpdateDeviceGroupPropertyByNameDefault describes a response with status code -1, with default header values.

Error
*/
type UpdateDeviceGroupPropertyByNameDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the update device group property by name default response
func (o *UpdateDeviceGroupPropertyByNameDefault) Code() int {
	return o._statusCode
}

func (o *UpdateDeviceGroupPropertyByNameDefault) Error() string {
	return fmt.Sprintf("[PUT /device/groups/{gid}/properties/{name}][%d] updateDeviceGroupPropertyByName default  %+v", o._statusCode, o.Payload)
}
func (o *UpdateDeviceGroupPropertyByNameDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateDeviceGroupPropertyByNameDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
