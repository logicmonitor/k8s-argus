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

// NewPatchReportGroupByIDParams creates a new PatchReportGroupByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPatchReportGroupByIDParams() *PatchReportGroupByIDParams {
	return &PatchReportGroupByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPatchReportGroupByIDParamsWithTimeout creates a new PatchReportGroupByIDParams object
// with the ability to set a timeout on a request.
func NewPatchReportGroupByIDParamsWithTimeout(timeout time.Duration) *PatchReportGroupByIDParams {
	return &PatchReportGroupByIDParams{
		timeout: timeout,
	}
}

// NewPatchReportGroupByIDParamsWithContext creates a new PatchReportGroupByIDParams object
// with the ability to set a context for a request.
func NewPatchReportGroupByIDParamsWithContext(ctx context.Context) *PatchReportGroupByIDParams {
	return &PatchReportGroupByIDParams{
		Context: ctx,
	}
}

// NewPatchReportGroupByIDParamsWithHTTPClient creates a new PatchReportGroupByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPatchReportGroupByIDParamsWithHTTPClient(client *http.Client) *PatchReportGroupByIDParams {
	return &PatchReportGroupByIDParams{
		HTTPClient: client,
	}
}

/* PatchReportGroupByIDParams contains all the parameters to send to the API endpoint
   for the patch report group by Id operation.

   Typically these are written to a http.Request.
*/
type PatchReportGroupByIDParams struct {

	// PatchFields.
	PatchFields *string

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Body.
	Body *models.ReportGroup

	// ID.
	//
	// Format: int32
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the patch report group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchReportGroupByIDParams) WithDefaults() *PatchReportGroupByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the patch report group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchReportGroupByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := PatchReportGroupByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithTimeout(timeout time.Duration) *PatchReportGroupByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithContext(ctx context.Context) *PatchReportGroupByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithHTTPClient(client *http.Client) *PatchReportGroupByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPatchFields adds the patchFields to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithPatchFields(patchFields *string) *PatchReportGroupByIDParams {
	o.SetPatchFields(patchFields)
	return o
}

// SetPatchFields adds the patchFields to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetPatchFields(patchFields *string) {
	o.PatchFields = patchFields
}

// WithUserAgent adds the userAgent to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithUserAgent(userAgent *string) *PatchReportGroupByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithBody adds the body to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithBody(body *models.ReportGroup) *PatchReportGroupByIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetBody(body *models.ReportGroup) {
	o.Body = body
}

// WithID adds the id to the patch report group by Id params
func (o *PatchReportGroupByIDParams) WithID(id int32) *PatchReportGroupByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch report group by Id params
func (o *PatchReportGroupByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchReportGroupByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
