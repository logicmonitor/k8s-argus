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

	"github.com/logicmonitor/lm-sdk-go/models"
)

// NewAddAlertNoteByIDParams creates a new AddAlertNoteByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewAddAlertNoteByIDParams() *AddAlertNoteByIDParams {
	return &AddAlertNoteByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewAddAlertNoteByIDParamsWithTimeout creates a new AddAlertNoteByIDParams object
// with the ability to set a timeout on a request.
func NewAddAlertNoteByIDParamsWithTimeout(timeout time.Duration) *AddAlertNoteByIDParams {
	return &AddAlertNoteByIDParams{
		timeout: timeout,
	}
}

// NewAddAlertNoteByIDParamsWithContext creates a new AddAlertNoteByIDParams object
// with the ability to set a context for a request.
func NewAddAlertNoteByIDParamsWithContext(ctx context.Context) *AddAlertNoteByIDParams {
	return &AddAlertNoteByIDParams{
		Context: ctx,
	}
}

// NewAddAlertNoteByIDParamsWithHTTPClient creates a new AddAlertNoteByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewAddAlertNoteByIDParamsWithHTTPClient(client *http.Client) *AddAlertNoteByIDParams {
	return &AddAlertNoteByIDParams{
		HTTPClient: client,
	}
}

/* AddAlertNoteByIDParams contains all the parameters to send to the API endpoint
   for the add alert note by Id operation.

   Typically these are written to a http.Request.
*/
type AddAlertNoteByIDParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Body.
	Body *models.AlertAck

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the add alert note by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddAlertNoteByIDParams) WithDefaults() *AddAlertNoteByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the add alert note by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddAlertNoteByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := AddAlertNoteByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithTimeout(timeout time.Duration) *AddAlertNoteByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithContext(ctx context.Context) *AddAlertNoteByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithHTTPClient(client *http.Client) *AddAlertNoteByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithUserAgent(userAgent *string) *AddAlertNoteByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithBody adds the body to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithBody(body *models.AlertAck) *AddAlertNoteByIDParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetBody(body *models.AlertAck) {
	o.Body = body
}

// WithID adds the id to the add alert note by Id params
func (o *AddAlertNoteByIDParams) WithID(id string) *AddAlertNoteByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the add alert note by Id params
func (o *AddAlertNoteByIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *AddAlertNoteByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
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
