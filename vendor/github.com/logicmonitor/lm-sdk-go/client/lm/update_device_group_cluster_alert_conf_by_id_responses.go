// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
	models "github.com/logicmonitor/lm-sdk-go/models"
)

// UpdateDeviceGroupClusterAlertConfByIDReader is a Reader for the UpdateDeviceGroupClusterAlertConfByID structure.
type UpdateDeviceGroupClusterAlertConfByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateDeviceGroupClusterAlertConfByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUpdateDeviceGroupClusterAlertConfByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewUpdateDeviceGroupClusterAlertConfByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateDeviceGroupClusterAlertConfByIDOK creates a UpdateDeviceGroupClusterAlertConfByIDOK with default headers values
func NewUpdateDeviceGroupClusterAlertConfByIDOK() *UpdateDeviceGroupClusterAlertConfByIDOK {
	return &UpdateDeviceGroupClusterAlertConfByIDOK{}
}

/*UpdateDeviceGroupClusterAlertConfByIDOK handles this case with default header values.

successful operation
*/
type UpdateDeviceGroupClusterAlertConfByIDOK struct {
	Payload *models.DeviceClusterAlertConfig
}

func (o *UpdateDeviceGroupClusterAlertConfByIDOK) Error() string {
	return fmt.Sprintf("[PUT /device/groups/{deviceGroupId}/clusterAlertConf/{id}][%d] updateDeviceGroupClusterAlertConfByIdOK  %+v", 200, o.Payload)
}

func (o *UpdateDeviceGroupClusterAlertConfByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.DeviceClusterAlertConfig)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateDeviceGroupClusterAlertConfByIDDefault creates a UpdateDeviceGroupClusterAlertConfByIDDefault with default headers values
func NewUpdateDeviceGroupClusterAlertConfByIDDefault(code int) *UpdateDeviceGroupClusterAlertConfByIDDefault {
	return &UpdateDeviceGroupClusterAlertConfByIDDefault{
		_statusCode: code,
	}
}

/*UpdateDeviceGroupClusterAlertConfByIDDefault handles this case with default header values.

Error
*/
type UpdateDeviceGroupClusterAlertConfByIDDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the update device group cluster alert conf by Id default response
func (o *UpdateDeviceGroupClusterAlertConfByIDDefault) Code() int {
	return o._statusCode
}

func (o *UpdateDeviceGroupClusterAlertConfByIDDefault) Error() string {
	return fmt.Sprintf("[PUT /device/groups/{deviceGroupId}/clusterAlertConf/{id}][%d] updateDeviceGroupClusterAlertConfById default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateDeviceGroupClusterAlertConfByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
