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

	"github.com/logicmonitor/lm-sdk-go/models"
)

// NewPatchDeviceGroupDatasourceByIDParams creates a new PatchDeviceGroupDatasourceByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPatchDeviceGroupDatasourceByIDParams() *PatchDeviceGroupDatasourceByIDParams {
	return &PatchDeviceGroupDatasourceByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPatchDeviceGroupDatasourceByIDParamsWithTimeout creates a new PatchDeviceGroupDatasourceByIDParams object
// with the ability to set a timeout on a request.
func NewPatchDeviceGroupDatasourceByIDParamsWithTimeout(timeout time.Duration) *PatchDeviceGroupDatasourceByIDParams {
	return &PatchDeviceGroupDatasourceByIDParams{
		timeout: timeout,
	}
}

// NewPatchDeviceGroupDatasourceByIDParamsWithContext creates a new PatchDeviceGroupDatasourceByIDParams object
// with the ability to set a context for a request.
func NewPatchDeviceGroupDatasourceByIDParamsWithContext(ctx context.Context) *PatchDeviceGroupDatasourceByIDParams {
	return &PatchDeviceGroupDatasourceByIDParams{
		Context: ctx,
	}
}

// NewPatchDeviceGroupDatasourceByIDParamsWithHTTPClient creates a new PatchDeviceGroupDatasourceByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPatchDeviceGroupDatasourceByIDParamsWithHTTPClient(client *http.Client) *PatchDeviceGroupDatasourceByIDParams {
	return &PatchDeviceGroupDatasourceByIDParams{
		HTTPClient: client,
	}
}

/* PatchDeviceGroupDatasourceByIDParams contains all the parameters to send to the API endpoint
   for the patch device group datasource by Id operation.

   Typically these are written to a http.Request.
*/
type PatchDeviceGroupDatasourceByIDParams struct {

	// PatchFields.
	PatchFields *string

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Body.
	Body *models.DeviceGroupDataSource

	// DeviceGroupID.
	//
	// Format: int32
	DeviceGroupID int32

	// ID.
	//
	// Format: int32
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the patch device group datasource by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchDeviceGroupDatasourceByIDParams) WithDefaults() *PatchDeviceGroupDatasourceByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the patch device group datasource by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchDeviceGroupDatasourceByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := PatchDeviceGroupDatasourceByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithTimeout(timeout time.Duration) *PatchDeviceGroupDatasourceByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithContext(ctx context.Context) *PatchDeviceGroupDatasourceByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithHTTPClient(client *http.Client) *PatchDeviceGroupDatasourceByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPatchFields adds the patchFields to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithPatchFields(patchFields *string) *PatchDeviceGroupDatasourceByIDParams {
	o.SetPatchFields(patchFields)
	return o
}

// SetPatchFields adds the patchFields to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetPatchFields(patchFields *string) {
	o.PatchFields = patchFields
}

// WithUserAgent adds the userAgent to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithUserAgent(userAgent *string) *PatchDeviceGroupDatasourceByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithBody adds the body to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithBody(body *models.DeviceGroupDataSource) *PatchDeviceGroupDatasourceByIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetBody(body *models.DeviceGroupDataSource) {
	o.Body = body
}

// WithDeviceGroupID adds the deviceGroupID to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithDeviceGroupID(deviceGroupID int32) *PatchDeviceGroupDatasourceByIDParams {
	o.SetDeviceGroupID(deviceGroupID)
	return o
}

// SetDeviceGroupID adds the deviceGroupId to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetDeviceGroupID(deviceGroupID int32) {
	o.DeviceGroupID = deviceGroupID
}

// WithID adds the id to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) WithID(id int32) *PatchDeviceGroupDatasourceByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch device group datasource by Id params
func (o *PatchDeviceGroupDatasourceByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchDeviceGroupDatasourceByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PatchFields != nil {

		// query param PatchFields
		var qrPatchFields string

		if o.PatchFields != nil {
			qrPatchFields = *o.PatchFields
		}
		qPatchFields := qrPatchFields
		if qPatchFields != "" {

			if err := r.SetQueryParam("PatchFields", qPatchFields); err != nil {
				return err
			}
		}
	}

	if o.UserAgent != nil {

		// header param User-Agent
		if err := r.SetHeaderParam("User-Agent", *o.UserAgent); err != nil {
			return err
		}
	}
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param deviceGroupId
	if err := r.SetPathParam("deviceGroupId", swag.FormatInt32(o.DeviceGroupID)); err != nil {
		return err
	}

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
