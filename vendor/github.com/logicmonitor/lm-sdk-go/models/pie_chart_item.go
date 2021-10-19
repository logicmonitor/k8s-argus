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

// PieChartItem pie chart item
//
// swagger:model PieChartItem
type PieChartItem struct {

	// color
	Color string `json:"color,omitempty"`

	// data point name
	// Required: true
	DataPointName *string `json:"dataPointName"`

	// legend
	// Required: true
	Legend *string `json:"legend"`
}

// Validate validates this pie chart item
func (m *PieChartItem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDataPointName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLegend(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PieChartItem) validateDataPointName(formats strfmt.Registry) error {

	if err := validate.Required("dataPointName", "body", m.DataPointName); err != nil {
		return err
	}

	return nil
}

func (m *PieChartItem) validateLegend(formats strfmt.Registry) error {

	if err := validate.Required("legend", "body", m.Legend); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this pie chart item based on context it is used
func (m *PieChartItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PieChartItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PieChartItem) UnmarshalBinary(b []byte) error {
	var res PieChartItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
