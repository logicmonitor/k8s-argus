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

// GetDeviceDatasourceListReader is a Reader for the GetDeviceDatasourceList structure.
type GetDeviceDatasourceListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceDatasourceListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDeviceDatasourceListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetDeviceDatasourceListTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetDeviceDatasourceListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDeviceDatasourceListOK creates a GetDeviceDatasourceListOK with default headers values
func NewGetDeviceDatasourceListOK() *GetDeviceDatasourceListOK {
	return &GetDeviceDatasourceListOK{}
}

/* GetDeviceDatasourceListOK describes a response with status code 200, with default header values.

successful operation
*/
type GetDeviceDatasourceListOK struct {
	Payload *models.DeviceDatasourcePaginationResponse
}

func (o *GetDeviceDatasourceListOK) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources][%d] getDeviceDatasourceListOK  %+v", 200, o.Payload)
}
func (o *GetDeviceDatasourceListOK) GetPayload() *models.DeviceDatasourcePaginationResponse {
	return o.Payload
}

func (o *GetDeviceDatasourceListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DeviceDatasourcePaginationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceDatasourceListTooManyRequests creates a GetDeviceDatasourceListTooManyRequests with default headers values
func NewGetDeviceDatasourceListTooManyRequests() *GetDeviceDatasourceListTooManyRequests {
	return &GetDeviceDatasourceListTooManyRequests{}
}

/* GetDeviceDatasourceListTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type GetDeviceDatasourceListTooManyRequests struct {

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

func (o *GetDeviceDatasourceListTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources][%d] getDeviceDatasourceListTooManyRequests ", 429)
}

func (o *GetDeviceDatasourceListTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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

// NewGetDeviceDatasourceListDefault creates a GetDeviceDatasourceListDefault with default headers values
func NewGetDeviceDatasourceListDefault(code int) *GetDeviceDatasourceListDefault {
	return &GetDeviceDatasourceListDefault{
		_statusCode: code,
	}
}

/* GetDeviceDatasourceListDefault describes a response with status code -1, with default header values.

Error
*/
type GetDeviceDatasourceListDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get device datasource list default response
func (o *GetDeviceDatasourceListDefault) Code() int {
	return o._statusCode
}

func (o *GetDeviceDatasourceListDefault) Error() string {
	return fmt.Sprintf("[GET /device/devices/{deviceId}/devicedatasources][%d] getDeviceDatasourceList default  %+v", o._statusCode, o.Payload)
}
func (o *GetDeviceDatasourceListDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetDeviceDatasourceListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
