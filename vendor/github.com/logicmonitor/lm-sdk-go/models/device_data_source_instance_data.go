// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DeviceDataSourceInstanceData device data source instance data
//
// swagger:model DeviceDataSourceInstanceData
type DeviceDataSourceInstanceData struct {

	// data source name
	// Read Only: true
	DataSourceName string `json:"dataSourceName,omitempty"`

	// next page params
	// Read Only: true
	NextPageParams string `json:"nextPageParams,omitempty"`

	// time
	// Read Only: true
	Time []int64 `json:"time,omitempty"`

	// values
	// Read Only: true
	Values [][]interface{} `json:"values,omitempty"`
}

// Validate validates this device data source instance data
func (m *DeviceDataSourceInstanceData) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this device data source instance data based on the context it is used
func (m *DeviceDataSourceInstanceData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDataSourceName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNextPageParams(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTime(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateValues(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DeviceDataSourceInstanceData) contextValidateDataSourceName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "dataSourceName", "body", string(m.DataSourceName)); err != nil {
		return err
	}

	return nil
}

func (m *DeviceDataSourceInstanceData) contextValidateNextPageParams(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "nextPageParams", "body", string(m.NextPageParams)); err != nil {
		return err
	}

	return nil
}

func (m *DeviceDataSourceInstanceData) contextValidateTime(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "time", "body", []int64(m.Time)); err != nil {
		return err
	}

	return nil
}

func (m *DeviceDataSourceInstanceData) contextValidateValues(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "values", "body", [][]interface{}(m.Values)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DeviceDataSourceInstanceData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceDataSourceInstanceData) UnmarshalBinary(b []byte) error {
	var res DeviceDataSourceInstanceData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
