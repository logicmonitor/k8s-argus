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

// NewAddNetscanParams creates a new AddNetscanParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewAddNetscanParams() *AddNetscanParams {
	return &AddNetscanParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewAddNetscanParamsWithTimeout creates a new AddNetscanParams object
// with the ability to set a timeout on a request.
func NewAddNetscanParamsWithTimeout(timeout time.Duration) *AddNetscanParams {
	return &AddNetscanParams{
		timeout: timeout,
	}
}

// NewAddNetscanParamsWithContext creates a new AddNetscanParams object
// with the ability to set a context for a request.
func NewAddNetscanParamsWithContext(ctx context.Context) *AddNetscanParams {
	return &AddNetscanParams{
		Context: ctx,
	}
}

// NewAddNetscanParamsWithHTTPClient creates a new AddNetscanParams object
// with the ability to set a custom HTTPClient for a request.
func NewAddNetscanParamsWithHTTPClient(client *http.Client) *AddNetscanParams {
	return &AddNetscanParams{
		HTTPClient: client,
	}
}

/* AddNetscanParams contains all the parameters to send to the API endpoint
   for the add netscan operation.

   Typically these are written to a http.Request.
*/
type AddNetscanParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Body.
	Body models.Netscan

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the add netscan params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddNetscanParams) WithDefaults() *AddNetscanParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the add netscan params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AddNetscanParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := AddNetscanParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the add netscan params
func (o *AddNetscanParams) WithTimeout(timeout time.Duration) *AddNetscanParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add netscan params
func (o *AddNetscanParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add netscan params
func (o *AddNetscanParams) WithContext(ctx context.Context) *AddNetscanParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add netscan params
func (o *AddNetscanParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add netscan params
func (o *AddNetscanParams) WithHTTPClient(client *http.Client) *AddNetscanParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add netscan params
func (o *AddNetscanParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the add netscan params
func (o *AddNetscanParams) WithUserAgent(userAgent *string) *AddNetscanParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the add netscan params
func (o *AddNetscanParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithBody adds the body to the add netscan params
func (o *AddNetscanParams) WithBody(body models.Netscan) *AddNetscanParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the add netscan params
func (o *AddNetscanParams) SetBody(body models.Netscan) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *AddNetscanParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
