// Code generated by go-swagger; DO NOT EDIT.

package lm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteDashboardByIDParams creates a new DeleteDashboardByIDParams object
// with the default values initialized.
func NewDeleteDashboardByIDParams() *DeleteDashboardByIDParams {
	var ()
	return &DeleteDashboardByIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteDashboardByIDParamsWithTimeout creates a new DeleteDashboardByIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteDashboardByIDParamsWithTimeout(timeout time.Duration) *DeleteDashboardByIDParams {
	var ()
	return &DeleteDashboardByIDParams{

		timeout: timeout,
	}
}

// NewDeleteDashboardByIDParamsWithContext creates a new DeleteDashboardByIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteDashboardByIDParamsWithContext(ctx context.Context) *DeleteDashboardByIDParams {
	var ()
	return &DeleteDashboardByIDParams{

		Context: ctx,
	}
}

// NewDeleteDashboardByIDParamsWithHTTPClient creates a new DeleteDashboardByIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteDashboardByIDParamsWithHTTPClient(client *http.Client) *DeleteDashboardByIDParams {
	var ()
	return &DeleteDashboardByIDParams{
		HTTPClient: client,
	}
}

/*DeleteDashboardByIDParams contains all the parameters to send to the API endpoint
for the delete dashboard by Id operation typically these are written to a http.Request
*/
type DeleteDashboardByIDParams struct {

	/*ID*/
	ID int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) WithTimeout(timeout time.Duration) *DeleteDashboardByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) WithContext(ctx context.Context) *DeleteDashboardByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) WithHTTPClient(client *http.Client) *DeleteDashboardByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) WithID(id int32) *DeleteDashboardByIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete dashboard by Id params
func (o *DeleteDashboardByIDParams) SetID(id int32) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteDashboardByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt32(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}