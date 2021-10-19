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
)

// StatsDGraph stats d graph
//
// swagger:model StatsDGraph
type StatsDGraph struct {

	// max value
	MaxValue float64 `json:"maxValue,omitempty"`

	// metrics
	Metrics []*StatsDMetricDefinition `json:"metrics,omitempty"`

	// min value
	MinValue float64 `json:"minValue,omitempty"`

	// namespace
	Namespace string `json:"namespace,omitempty"`

	// scale unit
	ScaleUnit int32 `json:"scaleUnit,omitempty"`

	// vertical label
	VerticalLabel string `json:"verticalLabel,omitempty"`
}

// Validate validates this stats d graph
func (m *StatsDGraph) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMetrics(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsDGraph) validateMetrics(formats strfmt.Registry) error {
	if swag.IsZero(m.Metrics) { // not required
		return nil
	}

	for i := 0; i < len(m.Metrics); i++ {
		if swag.IsZero(m.Metrics[i]) { // not required
			continue
		}

		if m.Metrics[i] != nil {
			if err := m.Metrics[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("metrics" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this stats d graph based on the context it is used
func (m *StatsDGraph) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMetrics(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsDGraph) contextValidateMetrics(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Metrics); i++ {

		if m.Metrics[i] != nil {
			if err := m.Metrics[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("metrics" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *StatsDGraph) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StatsDGraph) UnmarshalBinary(b []byte) error {
	var res StatsDGraph
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
