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

// AddDeviceReader is a Reader for the AddDevice structure.
type AddDeviceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddDeviceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAddDeviceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewAddDeviceTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewAddDeviceDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAddDeviceOK creates a AddDeviceOK with default headers values
func NewAddDeviceOK() *AddDeviceOK {
	return &AddDeviceOK{}
}

/* AddDeviceOK describes a response with status code 200, with default header values.

successful operation
*/
type AddDeviceOK struct {
	Payload *models.Device
}

func (o *AddDeviceOK) Error() string {
	return fmt.Sprintf("[POST /device/devices][%d] addDeviceOK  %+v", 200, o.Payload)
}
func (o *AddDeviceOK) GetPayload() *models.Device {
	return o.Payload
}

func (o *AddDeviceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Device)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddDeviceTooManyRequests creates a AddDeviceTooManyRequests with default headers values
func NewAddDeviceTooManyRequests() *AddDeviceTooManyRequests {
	return &AddDeviceTooManyRequests{}
}

/* AddDeviceTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type AddDeviceTooManyRequests struct {

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

func (o *AddDeviceTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /device/devices][%d] addDeviceTooManyRequests ", 429)
}

func (o *AddDeviceTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewAddDeviceDefault creates a AddDeviceDefault with default headers values
func NewAddDeviceDefault(code int) *AddDeviceDefault {
	return &AddDeviceDefault{
		_statusCode: code,
	}
}

/* AddDeviceDefault describes a response with status code -1, with default header values.

Error
*/
type AddDeviceDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the add device default response
func (o *AddDeviceDefault) Code() int {
	return o._statusCode
}

func (o *AddDeviceDefault) Error() string {
	return fmt.Sprintf("[POST /device/devices][%d] addDevice default  %+v", o._statusCode, o.Payload)
}
func (o *AddDeviceDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddDeviceDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
