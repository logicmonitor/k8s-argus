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

// NewDeleteReportGroupByIDParams creates a new DeleteReportGroupByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteReportGroupByIDParams() *DeleteReportGroupByIDParams {
	return &DeleteReportGroupByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteReportGroupByIDParamsWithTimeout creates a new DeleteReportGroupByIDParams object
// with the ability to set a timeout on a request.
func NewDeleteReportGroupByIDParamsWithTimeout(timeout time.Duration) *DeleteReportGroupByIDParams {
	return &DeleteReportGroupByIDParams{
		timeout: timeout,
	}
}

// NewDeleteReportGroupByIDParamsWithContext creates a new DeleteReportGroupByIDParams object
// with the ability to set a context for a request.
func NewDeleteReportGroupByIDParamsWithContext(ctx context.Context) *DeleteReportGroupByIDParams {
	return &DeleteReportGroupByIDParams{
		Context: ctx,
	}
}

// NewDeleteReportGroupByIDParamsWithHTTPClient creates a new DeleteReportGroupByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteReportGroupByIDParamsWithHTTPClient(client *http.Client) *DeleteReportGroupByIDParams {
	return &DeleteReportGroupByIDParams{
		HTTPClient: client,
	}
}

/* DeleteReportGroupByIDParams contains all the parameters to send to the API endpoint
   for the delete report group by Id operation.

   Typically these are written to a http.Request.
*/
type DeleteReportGroupByIDParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// ID.
	//
	// Format: int32
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete report group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteReportGroupByIDParams) WithDefaults() *DeleteReportGroupByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete report group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteReportGroupByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := DeleteReportGroupByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) WithTimeout(timeout time.Duration) *DeleteReportGroupByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) WithContext(ctx context.Context) *DeleteReportGroupByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) WithHTTPClient(client *http.Client) *DeleteReportGroupByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) WithUserAgent(userAgent *string) *DeleteReportGroupByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithID adds the id to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) WithID(id int32) *DeleteReportGroupByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete report group by Id params
func (o *DeleteReportGroupByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteReportGroupByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
