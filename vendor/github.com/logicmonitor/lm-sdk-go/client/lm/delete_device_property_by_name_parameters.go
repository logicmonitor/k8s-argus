// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewDeleteDevicePropertyByNameParams creates a new DeleteDevicePropertyByNameParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteDevicePropertyByNameParams() *DeleteDevicePropertyByNameParams {
	return &DeleteDevicePropertyByNameParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteDevicePropertyByNameParamsWithTimeout creates a new DeleteDevicePropertyByNameParams object
// with the ability to set a timeout on a request.
func NewDeleteDevicePropertyByNameParamsWithTimeout(timeout time.Duration) *DeleteDevicePropertyByNameParams {
	return &DeleteDevicePropertyByNameParams{
		timeout: timeout,
	}
}

// NewDeleteDevicePropertyByNameParamsWithContext creates a new DeleteDevicePropertyByNameParams object
// with the ability to set a context for a request.
func NewDeleteDevicePropertyByNameParamsWithContext(ctx context.Context) *DeleteDevicePropertyByNameParams {
	return &DeleteDevicePropertyByNameParams{
		Context: ctx,
	}
}

// NewDeleteDevicePropertyByNameParamsWithHTTPClient creates a new DeleteDevicePropertyByNameParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteDevicePropertyByNameParamsWithHTTPClient(client *http.Client) *DeleteDevicePropertyByNameParams {
	return &DeleteDevicePropertyByNameParams{
		HTTPClient: client,
	}
}

/* DeleteDevicePropertyByNameParams contains all the parameters to send to the API endpoint
   for the delete device property by name operation.

   Typically these are written to a http.Request.
*/
type DeleteDevicePropertyByNameParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// DeviceID.
	//
	// Format: int32
	DeviceID int32

	// Name.
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete device property by name params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteDevicePropertyByNameParams) WithDefaults() *DeleteDevicePropertyByNameParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete device property by name params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteDevicePropertyByNameParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := DeleteDevicePropertyByNameParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithTimeout(timeout time.Duration) *DeleteDevicePropertyByNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithContext(ctx context.Context) *DeleteDevicePropertyByNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithHTTPClient(client *http.Client) *DeleteDevicePropertyByNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithUserAgent(userAgent *string) *DeleteDevicePropertyByNameParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithDeviceID adds the deviceID to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithDeviceID(deviceID int32) *DeleteDevicePropertyByNameParams {
	o.SetDeviceID(deviceID)
	return o
}

// SetDeviceID adds the deviceId to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetDeviceID(deviceID int32) {
	o.DeviceID = deviceID
}

// WithName adds the name to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) WithName(name string) *DeleteDevicePropertyByNameParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the delete device property by name params
func (o *DeleteDevicePropertyByNameParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteDevicePropertyByNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.UserAgent != nil {

		// header param User-Agent
		if err := r.SetHeaderParam("User-Agent", *o.UserAgent); err != nil {
			return err
		}
	}

	// path param deviceId
	if err := r.SetPathParam("deviceId", swag.FormatInt32(o.DeviceID)); err != nil {
		return err
	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
