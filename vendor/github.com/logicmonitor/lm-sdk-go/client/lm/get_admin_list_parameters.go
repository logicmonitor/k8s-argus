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

// NewGetAdminListParams creates a new GetAdminListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetAdminListParams() *GetAdminListParams {
	return &GetAdminListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetAdminListParamsWithTimeout creates a new GetAdminListParams object
// with the ability to set a timeout on a request.
func NewGetAdminListParamsWithTimeout(timeout time.Duration) *GetAdminListParams {
	return &GetAdminListParams{
		timeout: timeout,
	}
}

// NewGetAdminListParamsWithContext creates a new GetAdminListParams object
// with the ability to set a context for a request.
func NewGetAdminListParamsWithContext(ctx context.Context) *GetAdminListParams {
	return &GetAdminListParams{
		Context: ctx,
	}
}

// NewGetAdminListParamsWithHTTPClient creates a new GetAdminListParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetAdminListParamsWithHTTPClient(client *http.Client) *GetAdminListParams {
	return &GetAdminListParams{
		HTTPClient: client,
	}
}

/* GetAdminListParams contains all the parameters to send to the API endpoint
   for the get admin list operation.

   Typically these are written to a http.Request.
*/
type GetAdminListParams struct {

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

// WithDefaults hydrates default values in the get admin list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAdminListParams) WithDefaults() *GetAdminListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get admin list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAdminListParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")

		offsetDefault = int32(0)

		sizeDefault = int32(50)
	)

	val := GetAdminListParams{
		UserAgent: &userAgentDefault,
		Offset:    &offsetDefault,
		Size:      &sizeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get admin list params
func (o *GetAdminListParams) WithTimeout(timeout time.Duration) *GetAdminListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get admin list params
func (o *GetAdminListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get admin list params
func (o *GetAdminListParams) WithContext(ctx context.Context) *GetAdminListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get admin list params
func (o *GetAdminListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get admin list params
func (o *GetAdminListParams) WithHTTPClient(client *http.Client) *GetAdminListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get admin list params
func (o *GetAdminListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get admin list params
func (o *GetAdminListParams) WithUserAgent(userAgent *string) *GetAdminListParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get admin list params
func (o *GetAdminListParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithFields adds the fields to the get admin list params
func (o *GetAdminListParams) WithFields(fields *string) *GetAdminListParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get admin list params
func (o *GetAdminListParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithFilter adds the filter to the get admin list params
func (o *GetAdminListParams) WithFilter(filter *string) *GetAdminListParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the get admin list params
func (o *GetAdminListParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithOffset adds the offset to the get admin list params
func (o *GetAdminListParams) WithOffset(offset *int32) *GetAdminListParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get admin list params
func (o *GetAdminListParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithSize adds the size to the get admin list params
func (o *GetAdminListParams) WithSize(size *int32) *GetAdminListParams {
	o.SetSize(size)
	return o
}

// SetSize adds the size to the get admin list params
func (o *GetAdminListParams) SetSize(size *int32) {
	o.Size = size
}

// WriteToRequest writes these params to a swagger request
func (o *GetAdminListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
