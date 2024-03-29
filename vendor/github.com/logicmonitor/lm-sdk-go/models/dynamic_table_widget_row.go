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

// DynamicTableWidgetRow dynamic table widget row
//
// swagger:model DynamicTableWidgetRow
type DynamicTableWidgetRow struct {

	// The display name of the device selected for the row
	// Read Only: true
	DeviceDisplayName string `json:"deviceDisplayName,omitempty"`

	// The full path of the group selected for the row
	// Read Only: true
	GroupFullPath string `json:"groupFullPath,omitempty"`

	// The instances for each column of the row
	// Read Only: true
	InstanceName string `json:"instanceName,omitempty"`

	// The label for the row
	Label string `json:"label,omitempty"`
}

// Validate validates this dynamic table widget row
func (m *DynamicTableWidgetRow) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this dynamic table widget row based on the context it is used
func (m *DynamicTableWidgetRow) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDeviceDisplayName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateGroupFullPath(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstanceName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DynamicTableWidgetRow) contextValidateDeviceDisplayName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "deviceDisplayName", "body", string(m.DeviceDisplayName)); err != nil {
		return err
	}

	return nil
}

func (m *DynamicTableWidgetRow) contextValidateGroupFullPath(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "groupFullPath", "body", string(m.GroupFullPath)); err != nil {
		return err
	}

	return nil
}

func (m *DynamicTableWidgetRow) contextValidateInstanceName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "instanceName", "body", string(m.InstanceName)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DynamicTableWidgetRow) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DynamicTableWidgetRow) UnmarshalBinary(b []byte) error {
	var res DynamicTableWidgetRow
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
