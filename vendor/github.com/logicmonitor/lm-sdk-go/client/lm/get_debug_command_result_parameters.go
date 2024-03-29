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

// NewGetDebugCommandResultParams creates a new GetDebugCommandResultParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDebugCommandResultParams() *GetDebugCommandResultParams {
	return &GetDebugCommandResultParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDebugCommandResultParamsWithTimeout creates a new GetDebugCommandResultParams object
// with the ability to set a timeout on a request.
func NewGetDebugCommandResultParamsWithTimeout(timeout time.Duration) *GetDebugCommandResultParams {
	return &GetDebugCommandResultParams{
		timeout: timeout,
	}
}

// NewGetDebugCommandResultParamsWithContext creates a new GetDebugCommandResultParams object
// with the ability to set a context for a request.
func NewGetDebugCommandResultParamsWithContext(ctx context.Context) *GetDebugCommandResultParams {
	return &GetDebugCommandResultParams{
		Context: ctx,
	}
}

// NewGetDebugCommandResultParamsWithHTTPClient creates a new GetDebugCommandResultParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDebugCommandResultParamsWithHTTPClient(client *http.Client) *GetDebugCommandResultParams {
	return &GetDebugCommandResultParams{
		HTTPClient: client,
	}
}

/* GetDebugCommandResultParams contains all the parameters to send to the API endpoint
   for the get debug command result operation.

   Typically these are written to a http.Request.
*/
type GetDebugCommandResultParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// CollectorID.
	//
	// Format: int32
	// Default: -1
	CollectorID *int32

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get debug command result params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDebugCommandResultParams) WithDefaults() *GetDebugCommandResultParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get debug command result params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDebugCommandResultParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")

		collectorIDDefault = int32(-1)
	)

	val := GetDebugCommandResultParams{
		UserAgent:   &userAgentDefault,
		CollectorID: &collectorIDDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get debug command result params
func (o *GetDebugCommandResultParams) WithTimeout(timeout time.Duration) *GetDebugCommandResultParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get debug command result params
func (o *GetDebugCommandResultParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get debug command result params
func (o *GetDebugCommandResultParams) WithContext(ctx context.Context) *GetDebugCommandResultParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get debug command result params
func (o *GetDebugCommandResultParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get debug command result params
func (o *GetDebugCommandResultParams) WithHTTPClient(client *http.Client) *GetDebugCommandResultParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get debug command result params
func (o *GetDebugCommandResultParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get debug command result params
func (o *GetDebugCommandResultParams) WithUserAgent(userAgent *string) *GetDebugCommandResultParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get debug command result params
func (o *GetDebugCommandResultParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithCollectorID adds the collectorID to the get debug command result params
func (o *GetDebugCommandResultParams) WithCollectorID(collectorID *int32) *GetDebugCommandResultParams {
	o.SetCollectorID(collectorID)
	return o
}

// SetCollectorID adds the collectorId to the get debug command result params
func (o *GetDebugCommandResultParams) SetCollectorID(collectorID *int32) {
	o.CollectorID = collectorID
}

// WithID adds the id to the get debug command result params
func (o *GetDebugCommandResultParams) WithID(id string) *GetDebugCommandResultParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get debug command result params
func (o *GetDebugCommandResultParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetDebugCommandResultParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.CollectorID != nil {

		// query param collectorId
		var qrCollectorID int32

		if o.CollectorID != nil {
			qrCollectorID = *o.CollectorID
		}
		qCollectorID := swag.FormatInt32(qrCollectorID)
		if qCollectorID != "" {

			if err := r.SetQueryParam("collectorId", qCollectorID); err != nil {
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
