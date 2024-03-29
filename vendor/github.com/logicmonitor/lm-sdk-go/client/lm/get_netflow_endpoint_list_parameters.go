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

// NewGetNetflowEndpointListParams creates a new GetNetflowEndpointListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetNetflowEndpointListParams() *GetNetflowEndpointListParams {
	return &GetNetflowEndpointListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetNetflowEndpointListParamsWithTimeout creates a new GetNetflowEndpointListParams object
// with the ability to set a timeout on a request.
func NewGetNetflowEndpointListParamsWithTimeout(timeout time.Duration) *GetNetflowEndpointListParams {
	return &GetNetflowEndpointListParams{
		timeout: timeout,
	}
}

// NewGetNetflowEndpointListParamsWithContext creates a new GetNetflowEndpointListParams object
// with the ability to set a context for a request.
func NewGetNetflowEndpointListParamsWithContext(ctx context.Context) *GetNetflowEndpointListParams {
	return &GetNetflowEndpointListParams{
		Context: ctx,
	}
}

// NewGetNetflowEndpointListParamsWithHTTPClient creates a new GetNetflowEndpointListParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetNetflowEndpointListParamsWithHTTPClient(client *http.Client) *GetNetflowEndpointListParams {
	return &GetNetflowEndpointListParams{
		HTTPClient: client,
	}
}

/* GetNetflowEndpointListParams contains all the parameters to send to the API endpoint
   for the get netflow endpoint list operation.

   Typically these are written to a http.Request.
*/
type GetNetflowEndpointListParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	// End.
	//
	// Format: int64
	End *int64

	// Fields.
	Fields *string

	// Filter.
	Filter *string

	// ID.
	//
	// Format: int32
	ID int32

	// NetflowFilter.
	NetflowFilter *string

	// Offset.
	//
	// Format: int32
	Offset *int32

	// Port.
	Port *string

	// Size.
	//
	// Format: int32
	// Default: 50
	Size *int32

	// Start.
	//
	// Format: int64
	Start *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get netflow endpoint list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetNetflowEndpointListParams) WithDefaults() *GetNetflowEndpointListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get netflow endpoint list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetNetflowEndpointListParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")

		offsetDefault = int32(0)

		sizeDefault = int32(50)
	)

	val := GetNetflowEndpointListParams{
		UserAgent: &userAgentDefault,
		Offset:    &offsetDefault,
		Size:      &sizeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithTimeout(timeout time.Duration) *GetNetflowEndpointListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithContext(ctx context.Context) *GetNetflowEndpointListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithHTTPClient(client *http.Client) *GetNetflowEndpointListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithUserAgent(userAgent *string) *GetNetflowEndpointListParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithEnd adds the end to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithEnd(end *int64) *GetNetflowEndpointListParams {
	o.SetEnd(end)
	return o
}

// SetEnd adds the end to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetEnd(end *int64) {
	o.End = end
}

// WithFields adds the fields to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithFields(fields *string) *GetNetflowEndpointListParams {
	o.SetFields(fields)
	return o
}

// SetFields adds the fields to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetFields(fields *string) {
	o.Fields = fields
}

// WithFilter adds the filter to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithFilter(filter *string) *GetNetflowEndpointListParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithID adds the id to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithID(id int32) *GetNetflowEndpointListParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetID(id int32) {
	o.ID = id
}

// WithNetflowFilter adds the netflowFilter to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithNetflowFilter(netflowFilter *string) *GetNetflowEndpointListParams {
	o.SetNetflowFilter(netflowFilter)
	return o
}

// SetNetflowFilter adds the netflowFilter to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetNetflowFilter(netflowFilter *string) {
	o.NetflowFilter = netflowFilter
}

// WithOffset adds the offset to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithOffset(offset *int32) *GetNetflowEndpointListParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetOffset(offset *int32) {
	o.Offset = offset
}

// WithPort adds the port to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithPort(port *string) *GetNetflowEndpointListParams {
	o.SetPort(port)
	return o
}

// SetPort adds the port to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetPort(port *string) {
	o.Port = port
}

// WithSize adds the size to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithSize(size *int32) *GetNetflowEndpointListParams {
	o.SetSize(size)
	return o
}

// SetSize adds the size to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetSize(size *int32) {
	o.Size = size
}

// WithStart adds the start to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) WithStart(start *int64) *GetNetflowEndpointListParams {
	o.SetStart(start)
	return o
}

// SetStart adds the start to the get netflow endpoint list params
func (o *GetNetflowEndpointListParams) SetStart(start *int64) {
	o.Start = start
}

// WriteToRequest writes these params to a swagger request
func (o *GetNetflowEndpointListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	if o.End != nil {

		// query param end
		var qrEnd int64

		if o.End != nil {
			qrEnd = *o.End
		}
		qEnd := swag.FormatInt64(qrEnd)
		if qEnd != "" {

			if err := r.SetQueryParam("end", qEnd); err != nil {
				return err
			}
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

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if o.NetflowFilter != nil {

		// query param netflowFilter
		var qrNetflowFilter string

		if o.NetflowFilter != nil {
			qrNetflowFilter = *o.NetflowFilter
		}
		qNetflowFilter := qrNetflowFilter
		if qNetflowFilter != "" {

			if err := r.SetQueryParam("netflowFilter", qNetflowFilter); err != nil {
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

	if o.Port != nil {

		// query param port
		var qrPort string

		if o.Port != nil {
			qrPort = *o.Port
		}
		qPort := qrPort
		if qPort != "" {

			if err := r.SetQueryParam("port", qPort); err != nil {
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

	if o.Start != nil {

		// query param start
		var qrStart int64

		if o.Start != nil {
			qrStart = *o.Start
		}
		qStart := swag.FormatInt64(qrStart)
		if qStart != "" {

			if err := r.SetQueryParam("start", qStart); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
