// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
	models "github.com/logicmonitor/lm-sdk-go/models"
	"golang.org/x/net/context"
)

// NewAddReportParams creates a new AddReportParams object
// with the default values initialized.
func NewAddReportParams() *AddReportParams {
	var ()
	return &AddReportParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAddReportParamsWithTimeout creates a new AddReportParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAddReportParamsWithTimeout(timeout time.Duration) *AddReportParams {
	var ()
	return &AddReportParams{

		timeout: timeout,
	}
}

// NewAddReportParamsWithContext creates a new AddReportParams object
// with the default values initialized, and the ability to set a context for a request
func NewAddReportParamsWithContext(ctx context.Context) *AddReportParams {
	var ()
	return &AddReportParams{

		Context: ctx,
	}
}

// NewAddReportParamsWithHTTPClient creates a new AddReportParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAddReportParamsWithHTTPClient(client *http.Client) *AddReportParams {
	var ()
	return &AddReportParams{
		HTTPClient: client,
	}
}

/*AddReportParams contains all the parameters to send to the API endpoint
for the add report operation typically these are written to a http.Request
*/
type AddReportParams struct {

	/*Body*/
	Body models.ReportBase

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the add report params
func (o *AddReportParams) WithTimeout(timeout time.Duration) *AddReportParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add report params
func (o *AddReportParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add report params
func (o *AddReportParams) WithContext(ctx context.Context) *AddReportParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add report params
func (o *AddReportParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add report params
func (o *AddReportParams) WithHTTPClient(client *http.Client) *AddReportParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add report params
func (o *AddReportParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the add report params
func (o *AddReportParams) WithBody(body models.ReportBase) *AddReportParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the add report params
func (o *AddReportParams) SetBody(body models.ReportBase) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *AddReportParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
