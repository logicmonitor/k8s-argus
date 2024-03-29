// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CellData cell data
//
// swagger:model CellData
type CellData struct {

	// alert severity
	AlertSeverity string `json:"alertSeverity,omitempty"`

	// alert status
	// Read Only: true
	AlertStatus string `json:"alertStatus,omitempty"`

	// days until alert list
	DaysUntilAlertList []*DaysUntilAlert `json:"daysUntilAlertList,omitempty"`

	// forecast day
	// Read Only: true
	ForecastDay int32 `json:"forecastDay,omitempty"`

	// instance Id
	// Read Only: true
	InstanceID int32 `json:"instanceId,omitempty"`

	// instance name
	// Read Only: true
	InstanceName string `json:"instanceName,omitempty"`

	// value
	// Read Only: true
	Value float64 `json:"value,omitempty"`
}

// Validate validates this cell data
func (m *CellData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDaysUntilAlertList(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CellData) validateDaysUntilAlertList(formats strfmt.Registry) error {
	if swag.IsZero(m.DaysUntilAlertList) { // not required
		return nil
	}

	for i := 0; i < len(m.DaysUntilAlertList); i++ {
		if swag.IsZero(m.DaysUntilAlertList[i]) { // not required
			continue
		}

		if m.DaysUntilAlertList[i] != nil {
			if err := m.DaysUntilAlertList[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("daysUntilAlertList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this cell data based on the context it is used
func (m *CellData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAlertStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDaysUntilAlertList(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateForecastDay(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstanceID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstanceName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateValue(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CellData) contextValidateAlertStatus(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "alertStatus", "body", string(m.AlertStatus)); err != nil {
		return err
	}

	return nil
}

func (m *CellData) contextValidateDaysUntilAlertList(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.DaysUntilAlertList); i++ {

		if m.DaysUntilAlertList[i] != nil {
			if err := m.DaysUntilAlertList[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("daysUntilAlertList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *CellData) contextValidateForecastDay(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "forecastDay", "body", int32(m.ForecastDay)); err != nil {
		return err
	}

	return nil
}

func (m *CellData) contextValidateInstanceID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "instanceId", "body", int32(m.InstanceID)); err != nil {
		return err
	}

	return nil
}

func (m *CellData) contextValidateInstanceName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "instanceName", "body", string(m.InstanceName)); err != nil {
		return err
	}

	return nil
}

func (m *CellData) contextValidateValue(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "value", "body", float64(m.Value)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CellData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CellData) UnmarshalBinary(b []byte) error {
	var res CellData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
