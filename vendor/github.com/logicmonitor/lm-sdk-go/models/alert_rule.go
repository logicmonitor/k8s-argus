// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AlertRule alert rule
// swagger:model AlertRule
type AlertRule struct {

	// The datapoint the alert rule is configured to match
	Datapoint string `json:"datapoint,omitempty"`

	// The datasource the alert rule is configured to match
	Datasource string `json:"datasource,omitempty"`

	// The device groups and service groups the alert rule is configured to match
	// Unique: true
	DeviceGroups []string `json:"deviceGroups,omitempty"`

	// The device names and service names the alert rule is configured to match
	// Unique: true
	Devices []string `json:"devices,omitempty"`

	// The escalation chain associated with the alert rule
	// Read Only: true
	EscalatingChain interface{} `json:"escalatingChain,omitempty"`

	// The id of the escalation chain associated with the alert rule
	// Required: true
	EscalatingChainID *int32 `json:"escalatingChainId"`

	// The escalation interval associated with the alert rule, in minutes
	EscalationInterval int32 `json:"escalationInterval,omitempty"`

	// The Id of the alert rule
	// Read Only: true
	ID int32 `json:"id,omitempty"`

	// The instance the alert rule is configured to match
	Instance string `json:"instance,omitempty"`

	// The alert severity levels the alert rule is configured to match. Acceptable values are: All, Warn, Error, Critical
	LevelStr string `json:"levelStr,omitempty"`

	// The name of the alert rule
	// Required: true
	Name *string `json:"name"`

	// The priority associated with the alert rule
	// Required: true
	Priority *int32 `json:"priority"`

	// Whether or not status notifications for acknowledgements and SDTs should be sent to the alert rule
	SuppressAlertAckSDT bool `json:"suppressAlertAckSdt,omitempty"`

	// Whether or not alert clear notifications should be sent to the alert rule
	SuppressAlertClear bool `json:"suppressAlertClear,omitempty"`
}

// Validate validates this alert rule
func (m *AlertRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDeviceGroups(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDevices(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEscalatingChainID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePriority(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertRule) validateDeviceGroups(formats strfmt.Registry) error {
	if swag.IsZero(m.DeviceGroups) { // not required
		return nil
	}

	if err := validate.UniqueItems("deviceGroups", "body", m.DeviceGroups); err != nil {
		return err
	}

	return nil
}

func (m *AlertRule) validateDevices(formats strfmt.Registry) error {
	if swag.IsZero(m.Devices) { // not required
		return nil
	}

	if err := validate.UniqueItems("devices", "body", m.Devices); err != nil {
		return err
	}

	return nil
}

func (m *AlertRule) validateEscalatingChainID(formats strfmt.Registry) error {
	if err := validate.Required("escalatingChainId", "body", m.EscalatingChainID); err != nil {
		return err
	}

	return nil
}

func (m *AlertRule) validateName(formats strfmt.Registry) error {
	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *AlertRule) validatePriority(formats strfmt.Registry) error {
	if err := validate.Required("priority", "body", m.Priority); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertRule) UnmarshalBinary(b []byte) error {
	var res AlertRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
