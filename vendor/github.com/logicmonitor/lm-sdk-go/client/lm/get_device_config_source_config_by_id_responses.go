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

// GetDeviceConfigSourceConfigByIDReader is a Reader for the GetDeviceConfigSourceConfigByID structure.
type GetDeviceConfigSourceConfigByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceConfigSourceConfigByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDeviceConfigSourceConfigByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetDeviceConfigSourceConfigByIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetDeviceConfigSourceConfigByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDeviceConfigSourceConfigByIDOK creates a GetDeviceConfigSourceConfigByIDOK with default headers values
func NewGetDeviceConfigSourceConfigByIDOK() *GetDeviceConfigSourceConfigByIDOK {
	return &GetDeviceConfigSourceConfigByIDOK{}
}

/* GetDeviceConfigSourceConfigByIDOK describes a response with status code 200, with default header values.

successful operation
*/
type GetDeviceConfigSourceConfigByIDOK struct {
	Payload *models.DeviceDataSourceInstanceConfig
}

func (o *GetDeviceConfigSourceConfigByIDOK) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{instanceId}/config/{id}][%d] getDeviceConfigSourceConfigByIdOK  %+v", 200, o.Payload)
}
func (o *GetDeviceConfigSourceConfigByIDOK) GetPayload() *models.DeviceDataSourceInstanceConfig {
	return o.Payload
}

func (o *GetDeviceConfigSourceConfigByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DeviceDataSourceInstanceConfig)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceConfigSourceConfigByIDTooManyRequests creates a GetDeviceConfigSourceConfigByIDTooManyRequests with default headers values
func NewGetDeviceConfigSourceConfigByIDTooManyRequests() *GetDeviceConfigSourceConfigByIDTooManyRequests {
	return &GetDeviceConfigSourceConfigByIDTooManyRequests{}
}

/* GetDeviceConfigSourceConfigByIDTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetDeviceConfigSourceConfigByIDTooManyRequests struct {

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

func (o *GetDeviceConfigSourceConfigByIDTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{instanceId}/config/{id}][%d] getDeviceConfigSourceConfigByIdTooManyRequests ", 429)
}

func (o *GetDeviceConfigSourceConfigByIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetDeviceConfigSourceConfigByIDDefault creates a GetDeviceConfigSourceConfigByIDDefault with default headers values
func NewGetDeviceConfigSourceConfigByIDDefault(code int) *GetDeviceConfigSourceConfigByIDDefault {
	return &GetDeviceConfigSourceConfigByIDDefault{
		_statusCode: code,
	}
}

/* GetDeviceConfigSourceConfigByIDDefault describes a response with status code -1, with default header values.

Error
*/
type GetDeviceConfigSourceConfigByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get device config source config by Id default response
func (o *GetDeviceConfigSourceConfigByIDDefault) Code() int {
	return o._statusCode
}

func (o *GetDeviceConfigSourceConfigByIDDefault) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources/{hdsId}/instances/{instanceId}/config/{id}][%d] getDeviceConfigSourceConfigById default  %+v", o._statusCode, o.Payload)
}
func (o *GetDeviceConfigSourceConfigByIDDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetDeviceConfigSourceConfigByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
