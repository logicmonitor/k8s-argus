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

// GetSiteMonitorCheckPointListReader is a Reader for the GetSiteMonitorCheckPointList structure.
type GetSiteMonitorCheckPointListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSiteMonitorCheckPointListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSiteMonitorCheckPointListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetSiteMonitorCheckPointListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetSiteMonitorCheckPointListOK creates a GetSiteMonitorCheckPointListOK with default headers values
func NewGetSiteMonitorCheckPointListOK() *GetSiteMonitorCheckPointListOK {
	return &GetSiteMonitorCheckPointListOK{}
}

/*GetSiteMonitorCheckPointListOK handles this case with default header values.

successful operation
*/
type GetSiteMonitorCheckPointListOK struct {
	Payload *models.SiteMonitorCheckPointPaginationResponse
}

func (o *GetSiteMonitorCheckPointListOK) Error() string {
	return fmt.Sprintf("[GET /website/smcheckpoints][%d] getSiteMonitorCheckPointListOK  %+v", 200, o.Payload)
}

func (o *GetSiteMonitorCheckPointListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.SiteMonitorCheckPointPaginationResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSiteMonitorCheckPointListDefault creates a GetSiteMonitorCheckPointListDefault with default headers values
func NewGetSiteMonitorCheckPointListDefault(code int) *GetSiteMonitorCheckPointListDefault {
	return &GetSiteMonitorCheckPointListDefault{
		_statusCode: code,
	}
}

/*GetSiteMonitorCheckPointListDefault handles this case with default header values.

Error
*/
type GetSiteMonitorCheckPointListDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get site monitor check point list default response
func (o *GetSiteMonitorCheckPointListDefault) Code() int {
	return o._statusCode
}

func (o *GetSiteMonitorCheckPointListDefault) Error() string {
	return fmt.Sprintf("[GET /website/smcheckpoints][%d] getSiteMonitorCheckPointList default  %+v", o._statusCode, o.Payload)
}

func (o *GetSiteMonitorCheckPointListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {
	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
