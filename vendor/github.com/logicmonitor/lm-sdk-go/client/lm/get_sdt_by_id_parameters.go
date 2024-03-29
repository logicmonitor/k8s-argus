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
)

// NewGetSDTByIDParams creates a new GetSDTByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetSDTByIDParams() *GetSDTByIDParams {
	return &GetSDTByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetSDTByIDParamsWithTimeout creates a new GetSDTByIDParams object
// with the ability to set a timeout on a request.
func NewGetSDTByIDParamsWithTimeout(timeout time.Duration) *GetSDTByIDParams {
	return &GetSDTByIDParams{
		timeout: timeout,
	}
}

// NewGetSDTByIDParamsWithContext creates a new GetSDTByIDParams object
// with the ability to set a context for a request.
func NewGetSDTByIDParamsWithContext(ctx context.Context) *GetSDTByIDParams {
	return &GetSDTByIDParams{
		Context: ctx,
	}
}

// NewGetSDTByIDParamsWithHTTPClient creates a new GetSDTByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetSDTByIDParamsWithHTTPClient(client *http.Client) *GetSDTByIDParams {
	return &GetSDTByIDParams{
		HTTPClient: client,
	}
}

/* GetSDTByIDParams contains all the parameters to send to the API endpoint
   for the get SDT by Id operation.

   Typically these are written to a http.Request.
*/
type GetSDTByIDParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Fields.
	Fields *string

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get SDT by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSDTByIDParams) WithDefaults() *GetSDTByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get SDT by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSDTByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := GetSDTByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get SDT by Id params
func (o *GetSDTByIDParams) WithTimeout(timeout time.Duration) *GetSDTByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get SDT by Id params
func (o *GetSDTByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get SDT by Id params
func (o *GetSDTByIDParams) WithContext(ctx context.Context) *GetSDTByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get SDT by Id params
func (o *GetSDTByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get SDT by Id params
func (o *GetSDTByIDParams) WithHTTPClient(client *http.Client) *GetSDTByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get SDT by Id params
func (o *GetSDTByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get SDT by Id params
func (o *GetSDTByIDParams) WithUserAgent(userAgent *string) *GetSDTByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get SDT by Id params
func (o *GetSDTByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithFields adds the fields to the get SDT by Id params
func (o *GetSDTByIDParams) WithFields(fields *string) *GetSDTByIDParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get SDT by Id params
func (o *GetSDTByIDParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithID adds the id to the get SDT by Id params
func (o *GetSDTByIDParams) WithID(id string) *GetSDTByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get SDT by Id params
func (o *GetSDTByIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetSDTByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
