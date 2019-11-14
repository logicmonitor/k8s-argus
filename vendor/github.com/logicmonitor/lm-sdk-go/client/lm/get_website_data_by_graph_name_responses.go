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

// GetWebsiteDataByGraphNameReader is a Reader for the GetWebsiteDataByGraphName structure.
type GetWebsiteDataByGraphNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWebsiteDataByGraphNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetWebsiteDataByGraphNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewGetWebsiteDataByGraphNameDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetWebsiteDataByGraphNameOK creates a GetWebsiteDataByGraphNameOK with default headers values
func NewGetWebsiteDataByGraphNameOK() *GetWebsiteDataByGraphNameOK {
	return &GetWebsiteDataByGraphNameOK{}
}

/*GetWebsiteDataByGraphNameOK handles this case with default header values.

successful operation
*/
type GetWebsiteDataByGraphNameOK struct {
	Payload *models.GraphPlot
}

func (o *GetWebsiteDataByGraphNameOK) Error() string {
	return fmt.Sprintf("[GET /website/websites/{id}/graphs/{graphName}/data][%d] getWebsiteDataByGraphNameOK  %+v", 200, o.Payload)
}

func (o *GetWebsiteDataByGraphNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GraphPlot)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetWebsiteDataByGraphNameDefault creates a GetWebsiteDataByGraphNameDefault with default headers values
func NewGetWebsiteDataByGraphNameDefault(code int) *GetWebsiteDataByGraphNameDefault {
	return &GetWebsiteDataByGraphNameDefault{
		_statusCode: code,
	}
}

/*GetWebsiteDataByGraphNameDefault handles this case with default header values.

Error
*/
type GetWebsiteDataByGraphNameDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get website data by graph name default response
func (o *GetWebsiteDataByGraphNameDefault) Code() int {
	return o._statusCode
}

func (o *GetWebsiteDataByGraphNameDefault) Error() string {
	return fmt.Sprintf("[GET /website/websites/{id}/graphs/{graphName}/data][%d] getWebsiteDataByGraphName default  %+v", o._statusCode, o.Payload)
}

func (o *GetWebsiteDataByGraphNameDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}