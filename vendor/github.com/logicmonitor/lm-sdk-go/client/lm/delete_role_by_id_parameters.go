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

// NewDeleteRoleByIDParams creates a new DeleteRoleByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteRoleByIDParams() *DeleteRoleByIDParams {
	return &DeleteRoleByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteRoleByIDParamsWithTimeout creates a new DeleteRoleByIDParams object
// with the ability to set a timeout on a request.
func NewDeleteRoleByIDParamsWithTimeout(timeout time.Duration) *DeleteRoleByIDParams {
	return &DeleteRoleByIDParams{
		timeout: timeout,
	}
}

// NewDeleteRoleByIDParamsWithContext creates a new DeleteRoleByIDParams object
// with the ability to set a context for a request.
func NewDeleteRoleByIDParamsWithContext(ctx context.Context) *DeleteRoleByIDParams {
	return &DeleteRoleByIDParams{
		Context: ctx,
	}
}

// NewDeleteRoleByIDParamsWithHTTPClient creates a new DeleteRoleByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteRoleByIDParamsWithHTTPClient(client *http.Client) *DeleteRoleByIDParams {
	return &DeleteRoleByIDParams{
		HTTPClient: client,
	}
}

/* DeleteRoleByIDParams contains all the parameters to send to the API endpoint
   for the delete role by Id operation.

   Typically these are written to a http.Request.
*/
type DeleteRoleByIDParams struct {

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

// WithDefaults hydrates default values in the delete role by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteRoleByIDParams) WithDefaults() *DeleteRoleByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete role by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteRoleByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := DeleteRoleByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the delete role by Id params
func (o *DeleteRoleByIDParams) WithTimeout(timeout time.Duration) *DeleteRoleByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete role by Id params
func (o *DeleteRoleByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete role by Id params
func (o *DeleteRoleByIDParams) WithContext(ctx context.Context) *DeleteRoleByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete role by Id params
func (o *DeleteRoleByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete role by Id params
func (o *DeleteRoleByIDParams) WithHTTPClient(client *http.Client) *DeleteRoleByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete role by Id params
func (o *DeleteRoleByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the delete role by Id params
func (o *DeleteRoleByIDParams) WithUserAgent(userAgent *string) *DeleteRoleByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the delete role by Id params
func (o *DeleteRoleByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithID adds the id to the delete role by Id params
func (o *DeleteRoleByIDParams) WithID(id int32) *DeleteRoleByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete role by Id params
func (o *DeleteRoleByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteRoleByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
