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

// AutomaticUpgradeInfo automatic upgrade info
//
// swagger:model AutomaticUpgradeInfo
type AutomaticUpgradeInfo struct {

	// created by
	// Read Only: true
	CreatedBy string `json:"createdBy,omitempty"`

	// day of week
	// Example: MON
	// Required: true
	DayOfWeek *string `json:"dayOfWeek"`

	// description
	// Example: regular MGD updates
	Description string `json:"description,omitempty"`

	// hour
	// Example: 15
	// Required: true
	Hour *int32 `json:"hour"`

	// level
	// Read Only: true
	Level string `json:"level,omitempty"`

	// minute
	// Example: 0
	// Required: true
	Minute *int32 `json:"minute"`

	// occurrence
	// Example: Any
	// Required: true
	Occurrence *string `json:"occurrence"`

	// timezone
	// Example: Americas/Los Angeles
	Timezone string `json:"timezone,omitempty"`

	// type
	// Read Only: true
	Type string `json:"type,omitempty"`

	// version
	// Example: MGD
	// Required: true
	Version *string `json:"version"`
}

// Validate validates this automatic upgrade info
func (m *AutomaticUpgradeInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDayOfWeek(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHour(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMinute(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOccurrence(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AutomaticUpgradeInfo) validateDayOfWeek(formats strfmt.Registry) error {

	if err := validate.Required("dayOfWeek", "body", m.DayOfWeek); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) validateHour(formats strfmt.Registry) error {

	if err := validate.Required("hour", "body", m.Hour); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) validateMinute(formats strfmt.Registry) error {

	if err := validate.Required("minute", "body", m.Minute); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) validateOccurrence(formats strfmt.Registry) error {

	if err := validate.Required("occurrence", "body", m.Occurrence); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this automatic upgrade info based on the context it is used
func (m *AutomaticUpgradeInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedBy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLevel(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AutomaticUpgradeInfo) contextValidateCreatedBy(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "createdBy", "body", string(m.CreatedBy)); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) contextValidateLevel(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "level", "body", string(m.Level)); err != nil {
		return err
	}

	return nil
}

func (m *AutomaticUpgradeInfo) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "type", "body", string(m.Type)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AutomaticUpgradeInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AutomaticUpgradeInfo) UnmarshalBinary(b []byte) error {
	var res AutomaticUpgradeInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
