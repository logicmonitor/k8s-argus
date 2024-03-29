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
)

// NewImportDNSMappingParams creates a new ImportDNSMappingParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportDNSMappingParams() *ImportDNSMappingParams {
	return &ImportDNSMappingParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportDNSMappingParamsWithTimeout creates a new ImportDNSMappingParams object
// with the ability to set a timeout on a request.
func NewImportDNSMappingParamsWithTimeout(timeout time.Duration) *ImportDNSMappingParams {
	return &ImportDNSMappingParams{
		timeout: timeout,
	}
}

// NewImportDNSMappingParamsWithContext creates a new ImportDNSMappingParams object
// with the ability to set a context for a request.
func NewImportDNSMappingParamsWithContext(ctx context.Context) *ImportDNSMappingParams {
	return &ImportDNSMappingParams{
		Context: ctx,
	}
}

// NewImportDNSMappingParamsWithHTTPClient creates a new ImportDNSMappingParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportDNSMappingParamsWithHTTPClient(client *http.Client) *ImportDNSMappingParams {
	return &ImportDNSMappingParams{
		HTTPClient: client,
	}
}

/* ImportDNSMappingParams contains all the parameters to send to the API endpoint
   for the import DNS mapping operation.

   Typically these are written to a http.Request.
*/
type ImportDNSMappingParams struct {

	// UserAgent.
	//
	// Default: "Logicmonitor/SDK: Argus Dist-v1.0.0-argus1"
	UserAgent *string

	/* File.

	   the csv mapping to be uploaded
	*/
	File runtime.NamedReadCloser

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import DNS mapping params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportDNSMappingParams) WithDefaults() *ImportDNSMappingParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import DNS mapping params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportDNSMappingParams) SetDefaults() {
	var (
		userAgentDefault = string("Logicmonitor/SDK: Argus Dist-v1.0.0-argus1")
	)

	val := ImportDNSMappingParams{
		UserAgent: &userAgentDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the import DNS mapping params
func (o *ImportDNSMappingParams) WithTimeout(timeout time.Duration) *ImportDNSMappingParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import DNS mapping params
func (o *ImportDNSMappingParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import DNS mapping params
func (o *ImportDNSMappingParams) WithContext(ctx context.Context) *ImportDNSMappingParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import DNS mapping params
func (o *ImportDNSMappingParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import DNS mapping params
func (o *ImportDNSMappingParams) WithHTTPClient(client *http.Client) *ImportDNSMappingParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import DNS mapping params
func (o *ImportDNSMappingParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithUserAgent adds the userAgent to the import DNS mapping params
func (o *ImportDNSMappingParams) WithUserAgent(userAgent *string) *ImportDNSMappingParams {
	o.SetUserAgent(userAgent)
	return o
}

// SetUserAgent adds the userAgent to the import DNS mapping params
func (o *ImportDNSMappingParams) SetUserAgent(userAgent *string) {
	o.UserAgent = userAgent
}

// WithFile adds the file to the import DNS mapping params
func (o *ImportDNSMappingParams) WithFile(file runtime.NamedReadCloser) *ImportDNSMappingParams {
	o.SetFile(file)
	return o
}

// SetFile adds the file to the import DNS mapping params
func (o *ImportDNSMappingParams) SetFile(file runtime.NamedReadCloser) {
	o.File = file
}

// WriteToRequest writes these params to a swagger request
func (o *ImportDNSMappingParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
	// form file param file
	if err := r.SetFileParam("file", o.File); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
