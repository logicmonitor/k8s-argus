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

// NewGetWidgetByIDParams creates a new GetWidgetByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetWidgetByIDParams() *GetWidgetByIDParams {
	return &GetWidgetByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetWidgetByIDParamsWithTimeout creates a new GetWidgetByIDParams object
// with the ability to set a timeout on a request.
func NewGetWidgetByIDParamsWithTimeout(timeout time.Duration) *GetWidgetByIDParams {
	return &GetWidgetByIDParams{
		timeout: timeout,
	}
}

// NewGetWidgetByIDParamsWithContext creates a new GetWidgetByIDParams object
// with the ability to set a context for a request.
func NewGetWidgetByIDParamsWithContext(ctx context.Context) *GetWidgetByIDParams {
	return &GetWidgetByIDParams{
		Context: ctx,
	}
}

// NewGetWidgetByIDParamsWithHTTPClient creates a new GetWidgetByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetWidgetByIDParamsWithHTTPClient(client *http.Client) *GetWidgetByIDParams {
	return &GetWidgetByIDParams{
		HTTPClient: client,
	}
}

/* GetWidgetByIDParams contains all the parameters to send to the API endpoint
   for the get widget by Id operation.

   Typically these are written to a http.Request.
*/
type GetWidgetByIDParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Fields.
	Fields *string

	// ID.
	//
	// Format: int32
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get widget by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetWidgetByIDParams) WithDefaults() *GetWidgetByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get widget by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetWidgetByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := GetWidgetByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get widget by Id params
func (o *GetWidgetByIDParams) WithTimeout(timeout time.Duration) *GetWidgetByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get widget by Id params
func (o *GetWidgetByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get widget by Id params
func (o *GetWidgetByIDParams) WithContext(ctx context.Context) *GetWidgetByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get widget by Id params
func (o *GetWidgetByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get widget by Id params
func (o *GetWidgetByIDParams) WithHTTPClient(client *http.Client) *GetWidgetByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get widget by Id params
func (o *GetWidgetByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get widget by Id params
func (o *GetWidgetByIDParams) WithUserAgent(userAgent *string) *GetWidgetByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get widget by Id params
func (o *GetWidgetByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithFields adds the fields to the get widget by Id params
func (o *GetWidgetByIDParams) WithFields(fields *string) *GetWidgetByIDParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get widget by Id params
func (o *GetWidgetByIDParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithID adds the id to the get widget by Id params
func (o *GetWidgetByIDParams) WithID(id int32) *GetWidgetByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get widget by Id params
func (o *GetWidgetByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetWidgetByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
