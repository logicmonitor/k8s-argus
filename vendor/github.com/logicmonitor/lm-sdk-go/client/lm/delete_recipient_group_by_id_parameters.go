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

// NewDeleteRecipientGroupByIDParams creates a new DeleteRecipientGroupByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteRecipientGroupByIDParams() *DeleteRecipientGroupByIDParams {
	return &DeleteRecipientGroupByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteRecipientGroupByIDParamsWithTimeout creates a new DeleteRecipientGroupByIDParams object
// with the ability to set a timeout on a request.
func NewDeleteRecipientGroupByIDParamsWithTimeout(timeout time.Duration) *DeleteRecipientGroupByIDParams {
	return &DeleteRecipientGroupByIDParams{
		timeout: timeout,
	}
}

// NewDeleteRecipientGroupByIDParamsWithContext creates a new DeleteRecipientGroupByIDParams object
// with the ability to set a context for a request.
func NewDeleteRecipientGroupByIDParamsWithContext(ctx context.Context) *DeleteRecipientGroupByIDParams {
	return &DeleteRecipientGroupByIDParams{
		Context: ctx,
	}
}

// NewDeleteRecipientGroupByIDParamsWithHTTPClient creates a new DeleteRecipientGroupByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteRecipientGroupByIDParamsWithHTTPClient(client *http.Client) *DeleteRecipientGroupByIDParams {
	return &DeleteRecipientGroupByIDParams{
		HTTPClient: client,
	}
}

/* DeleteRecipientGroupByIDParams contains all the parameters to send to the API endpoint
   for the delete recipient group by Id operation.

   Typically these are written to a http.Request.
*/
type DeleteRecipientGroupByIDParams struct {

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

// WithDefaults hydrates default values in the delete recipient group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteRecipientGroupByIDParams) WithDefaults() *DeleteRecipientGroupByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete recipient group by Id params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteRecipientGroupByIDParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := DeleteRecipientGroupByIDParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) WithTimeout(timeout time.Duration) *DeleteRecipientGroupByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) WithContext(ctx context.Context) *DeleteRecipientGroupByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) WithHTTPClient(client *http.Client) *DeleteRecipientGroupByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) WithUserAgent(userAgent *string) *DeleteRecipientGroupByIDParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithID adds the id to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) WithID(id int32) *DeleteRecipientGroupByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete recipient group by Id params
func (o *DeleteRecipientGroupByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteRecipientGroupByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
