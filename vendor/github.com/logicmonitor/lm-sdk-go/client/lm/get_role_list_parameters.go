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

// NewGetRoleListParams creates a new GetRoleListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetRoleListParams() *GetRoleListParams {
	return &GetRoleListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetRoleListParamsWithTimeout creates a new GetRoleListParams object
// with the ability to set a timeout on a request.
func NewGetRoleListParamsWithTimeout(timeout time.Duration) *GetRoleListParams {
	return &GetRoleListParams{
		timeout: timeout,
	}
}

// NewGetRoleListParamsWithContext creates a new GetRoleListParams object
// with the ability to set a context for a request.
func NewGetRoleListParamsWithContext(ctx context.Context) *GetRoleListParams {
	return &GetRoleListParams{
		Context: ctx,
	}
}

// NewGetRoleListParamsWithHTTPClient creates a new GetRoleListParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetRoleListParamsWithHTTPClient(client *http.Client) *GetRoleListParams {
	return &GetRoleListParams{
		HTTPClient: client,
	}
}

/* GetRoleListParams contains all the parameters to send to the API endpoint
   for the get role list operation.

   Typically these are written to a http.Request.
*/
type GetRoleListParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// Fields.
	Fields *string

	// Filter.
	Filter *string

	// Offset.
	//
	// Format: int32
	Offset *int32

	// Size.
	//
	// Format: int32
	// Default: 50
	Size *int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get role list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRoleListParams) WithDefaults() *GetRoleListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get role list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRoleListParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")

		offsetDefault = int32(0)

		sizeDefault = int32(50)
	)

	val := GetRoleListParams{
		UserAgent: &userAgentDefault,
		Offset:    &offsetDefault,
		Size:      &sizeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get role list params
func (o *GetRoleListParams) WithTimeout(timeout time.Duration) *GetRoleListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get role list params
func (o *GetRoleListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get role list params
func (o *GetRoleListParams) WithContext(ctx context.Context) *GetRoleListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get role list params
func (o *GetRoleListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get role list params
func (o *GetRoleListParams) WithHTTPClient(client *http.Client) *GetRoleListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get role list params
func (o *GetRoleListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get role list params
func (o *GetRoleListParams) WithUserAgent(userAgent *string) *GetRoleListParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get role list params
func (o *GetRoleListParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithFields adds the fields to the get role list params
func (o *GetRoleListParams) WithFields(fields *string) *GetRoleListParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get role list params
func (o *GetRoleListParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithFilter adds the filter to the get role list params
func (o *GetRoleListParams) WithFilter(filter *string) *GetRoleListParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the get role list params
func (o *GetRoleListParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithOffset adds the offset to the get role list params
func (o *GetRoleListParams) WithOffset(offset *int32) *GetRoleListParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get role list params
func (o *GetRoleListParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithSize adds the size to the get role list params
func (o *GetRoleListParams) WithSize(size *int32) *GetRoleListParams {
	o.SetSize(size)
	return o
}

// SetSize adds the size to the get role list params
func (o *GetRoleListParams) SetSize(size *int32) {
	o.Size = size
}

// WriteToRequest writes these params to a swagger request
func (o *GetRoleListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.Filter != nil {

		// query param filter
		var qrFilter string

		if o.Filter != nil {
			qrFilter = *o.Filter
		}
		qFilter := qrFilter
		if qFilter != "" {

			if err := r.SetQueryParam("filter", qFilter); err != nil {
				return err
			}
		}
	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int32

		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt32(qrOffset)
		if qOffset != "" {

			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}
	}

	if o.Size != nil {

		// query param size
		var qrSize int32

		if o.Size != nil {
			qrSize = *o.Size
		}
		qSize := swag.FormatInt32(qrSize)
		if qSize != "" {

			if err := r.SetQueryParam("size", qSize); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
