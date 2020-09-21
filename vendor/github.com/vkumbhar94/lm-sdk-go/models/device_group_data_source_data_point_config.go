// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DeviceGroupDataSourceDataPointConfig device group data source data point config
// swagger:model DeviceGroupDataSourceDataPointConfig
type DeviceGroupDataSourceDataPointConfig struct {

	// alert expr
	// Required: true
	AlertExpr *string `json:"alertExpr"`

	// alert expr note
	AlertExprNote string `json:"alertExprNote,omitempty"`

	// data point description
	// Read Only: true
	DataPointDescription string `json:"dataPointDescription,omitempty"`

	// data point Id
	// Required: true
	DataPointID *int32 `json:"dataPointId"`

	// data point name
	// Required: true
	DataPointName *string `json:"dataPointName"`

	// disable alerting
	DisableAlerting bool `json:"disableAlerting,omitempty"`

	// global alert expr
	// Read Only: true
	GlobalAlertExpr string `json:"globalAlertExpr,omitempty"`
}

// Validate validates this device group data source data point config
func (m *DeviceGroupDataSourceDataPointConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlertExpr(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDataPointID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDataPointName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DeviceGroupDataSourceDataPointConfig) validateAlertExpr(formats strfmt.Registry) error {

	if err := validate.Required("alertExpr", "body", m.AlertExpr); err != nil {
		return err
	}

	return nil
}

func (m *DeviceGroupDataSourceDataPointConfig) validateDataPointID(formats strfmt.Registry) error {

	if err := validate.Required("dataPointId", "body", m.DataPointID); err != nil {
		return err
	}

	return nil
}

func (m *DeviceGroupDataSourceDataPointConfig) validateDataPointName(formats strfmt.Registry) error {

	if err := validate.Required("dataPointName", "body", m.DataPointName); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DeviceGroupDataSourceDataPointConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceGroupDataSourceDataPointConfig) UnmarshalBinary(b []byte) error {
	var res DeviceGroupDataSourceDataPointConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}