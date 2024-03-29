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

// NewPatchRoleByIDParams creates a new PatchRoleByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPatchRoleByIDParams() *PatchRoleByIDParams {
	return &PatchRoleByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPatchRoleByIDParamsWithTimeout creates a new PatchRoleByIDParams object
// with the ability to set a timeout on a request.
func NewPatchRoleByIDParamsWithTimeout(timeout time.Duration) *PatchRoleByIDParams {
	return &PatchRoleByIDParams{
		timeout: timeout,
	}
}

// NewPatchRoleByIDParamsWithContext creates a new PatchRoleByIDParams object
// with the ability to set a context for a request.
func NewPatchRoleByIDParamsWithContext(ctx context.Context) *PatchRoleByIDParams {
	return &PatchRoleByIDParams{
		Context: ctx,
	}
}

// NewPatchRoleByIDParamsWithHTTPClient creates a new PatchRoleByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewPatchRoleByIDParamsWithHTTPClient(client *http.Client) *PatchRoleByIDParams {
	return &PatchRoleByIDParams{
		HTTPClient: client,
	}
}

/* PatchRoleByIDParams contains all the parameters to send to the API endpoint
   for the patch role by Id operation.

   Typically these are written to a http.Request.
*/
type PatchRoleByIDParams struct {

	// PatchFields.
	PatchFields *string

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Body.
	Body *models.Role

	// ID.
	//
	// Format: int32
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the patch role by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchRoleByIDParams) WithDefaults() *PatchRoleByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the patch role by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PatchRoleByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := PatchRoleByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the patch role by Id params
func (o *PatchRoleByIDParams) WithTimeout(timeout time.Duration) *PatchRoleByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch role by Id params
func (o *PatchRoleByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch role by Id params
func (o *PatchRoleByIDParams) WithContext(ctx context.Context) *PatchRoleByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch role by Id params
func (o *PatchRoleByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch role by Id params
func (o *PatchRoleByIDParams) WithHTTPClient(client *http.Client) *PatchRoleByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch role by Id params
func (o *PatchRoleByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPatchFields adds the patchFields to the patch role by Id params
func (o *PatchRoleByIDParams) WithPatchFields(patchFields *string) *PatchRoleByIDParams {
	o.SetPatchFields(patchFields)
	return o
}

// SetPatchFields adds the patchFields to the patch role by Id params
func (o *PatchRoleByIDParams) SetPatchFields(patchFields *string) {
	o.PatchFields = patchFields
}

// WithUserAgent adds the userAgent to the patch role by Id params
func (o *PatchRoleByIDParams) WithUserAgent(userAgent *string) *PatchRoleByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the patch role by Id params
func (o *PatchRoleByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithBody adds the body to the patch role by Id params
func (o *PatchRoleByIDParams) WithBody(body *models.Role) *PatchRoleByIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the patch role by Id params
func (o *PatchRoleByIDParams) SetBody(body *models.Role) {
	o.Body = body
}

// WithID adds the id to the patch role by Id params
func (o *PatchRoleByIDParams) WithID(id int32) *PatchRoleByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch role by Id params
func (o *PatchRoleByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchRoleByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
