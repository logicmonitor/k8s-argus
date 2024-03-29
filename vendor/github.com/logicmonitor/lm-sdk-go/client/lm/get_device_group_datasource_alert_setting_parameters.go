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

// NewGetDeviceGroupDatasourceAlertSettingParams creates a new GetDeviceGroupDatasourceAlertSettingParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDeviceGroupDatasourceAlertSettingParams() *GetDeviceGroupDatasourceAlertSettingParams {
	return &GetDeviceGroupDatasourceAlertSettingParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDeviceGroupDatasourceAlertSettingParamsWithTimeout creates a new GetDeviceGroupDatasourceAlertSettingParams object
// with the ability to set a timeout on a request.
func NewGetDeviceGroupDatasourceAlertSettingParamsWithTimeout(timeout time.Duration) *GetDeviceGroupDatasourceAlertSettingParams {
	return &GetDeviceGroupDatasourceAlertSettingParams{
		timeout: timeout,
	}
}

// NewGetDeviceGroupDatasourceAlertSettingParamsWithContext creates a new GetDeviceGroupDatasourceAlertSettingParams object
// with the ability to set a context for a request.
func NewGetDeviceGroupDatasourceAlertSettingParamsWithContext(ctx context.Context) *GetDeviceGroupDatasourceAlertSettingParams {
	return &GetDeviceGroupDatasourceAlertSettingParams{
		Context: ctx,
	}
}

// NewGetDeviceGroupDatasourceAlertSettingParamsWithHTTPClient creates a new GetDeviceGroupDatasourceAlertSettingParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDeviceGroupDatasourceAlertSettingParamsWithHTTPClient(client *http.Client) *GetDeviceGroupDatasourceAlertSettingParams {
	return &GetDeviceGroupDatasourceAlertSettingParams{
		HTTPClient: client,
	}
}

/* GetDeviceGroupDatasourceAlertSettingParams contains all the parameters to send to the API endpoint
   for the get device group datasource alert setting operation.

   Typically these are written to a http.Request.
*/
type GetDeviceGroupDatasourceAlertSettingParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// DeviceGroupID.
	//
	// Format: int32
	DeviceGroupID int32

	// DsID.
	//
	// Format: int32
	DsID int32

	// Fields.
	Fields *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get device group datasource alert setting params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithDefaults() *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get device group datasource alert setting params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := GetDeviceGroupDatasourceAlertSettingParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithTimeout(timeout time.Duration) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithContext(ctx context.Context) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithHTTPClient(client *http.Client) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithUserAgent(userAgent *string) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithDeviceGroupID adds the deviceGroupID to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithDeviceGroupID(deviceGroupID int32) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetDeviceGroupID(deviceGroupID)
	return o
}

// SetDeviceGroupID adds the deviceGroupId to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetDeviceGroupID(deviceGroupID int32) {
	o.DeviceGroupID = deviceGroupID
}

// WithDsID adds the dsID to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithDsID(dsID int32) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetDsID(dsID)
	return o
}

// SetDsID adds the dsId to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetDsID(dsID int32) {
	o.DsID = dsID
}

// WithFields adds the fields to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) WithFields(fields *string) *GetDeviceGroupDatasourceAlertSettingParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get device group datasource alert setting params
func (o *GetDeviceGroupDatasourceAlertSettingParams) SetFields(fields *string) {
	o.Fields = fields
}

// WriteToRequest writes these params to a swagger request
func (o *GetDeviceGroupDatasourceAlertSettingParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param deviceGroupId
	if err := r.SetPathParam("deviceGroupId", swag.FormatInt32(o.DeviceGroupID)); err != nil {
		return err
	}

	// path param dsId
	if err := r.SetPathParam("dsId", swag.FormatInt32(o.DsID)); err != nil {
		return err
	}

	if o.Fields != nil {

		// query param fields
		var qrFields string

		if o.Fields != nil {
			qrFields = *o.Fields
		}
		qFields := qrFields
		if qFields != "" {

			if err := r.SetQueryParam("fields", qFields); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
