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

// PieChartInfo pie chart info
//
// swagger:model PieChartInfo
type PieChartInfo struct {

	// The counter is used for saving applyTo expression, it's mainly used for count device
	Counters []*Counter `json:"counters,omitempty"`

	// The datapoints added to the widget. Note that datapoints must be included in the pieChartItems object to be displayed in the widget
	DataPoints []*PieChartDataPoint `json:"dataPoints,omitempty"`

	// If the number of slices exceeds the maxSlicesCanBeShown, this value indicates whether the remaining slices should be grouped together
	GroupRemainingAsOthers bool `json:"groupRemainingAsOthers,omitempty"`

	// Whether items at 0% should be hidden
	HideZeroPercentSlices bool `json:"hideZeroPercentSlices,omitempty"`

	// The maximum number of slices you'd like displayed in the pie chart
	MaxSlicesCanBeShown int32 `json:"maxSlicesCanBeShown,omitempty"`

	// The datapoints and virtual datapoints that will be displayed in the pie chart
	// Required: true
	PieChartItems []*PieChartItem `json:"pieChartItems"`

	// Whether or not labels and lines should be displayed on the pie chart
	ShowLabelsAndLinesOnPC bool `json:"showLabelsAndLinesOnPC,omitempty"`

	// The title that will be displayed above the pie chart
	Title string `json:"title,omitempty"`

	// The virtual datapoints added to the widget. Note that virtual datapoints must be included in the pieChartItems object to be displayed in the widget
	VirtualDataPoints []*VirtualDataPoint `json:"virtualDataPoints,omitempty"`
}

// Validate validates this pie chart info
func (m *PieChartInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCounters(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDataPoints(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePieChartItems(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVirtualDataPoints(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PieChartInfo) validateCounters(formats strfmt.Registry) error {
	if swag.IsZero(m.Counters) { // not required
		return nil
	}

	for i := 0; i < len(m.Counters); i++ {
		if swag.IsZero(m.Counters[i]) { // not required
			continue
		}

		if m.Counters[i] != nil {
			if err := m.Counters[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("counters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) validateDataPoints(formats strfmt.Registry) error {
	if swag.IsZero(m.DataPoints) { // not required
		return nil
	}

	for i := 0; i < len(m.DataPoints); i++ {
		if swag.IsZero(m.DataPoints[i]) { // not required
			continue
		}

		if m.DataPoints[i] != nil {
			if err := m.DataPoints[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("dataPoints" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) validatePieChartItems(formats strfmt.Registry) error {

	if err := validate.Required("pieChartItems", "body", m.PieChartItems); err != nil {
		return err
	}

	for i := 0; i < len(m.PieChartItems); i++ {
		if swag.IsZero(m.PieChartItems[i]) { // not required
			continue
		}

		if m.PieChartItems[i] != nil {
			if err := m.PieChartItems[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("pieChartItems" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) validateVirtualDataPoints(formats strfmt.Registry) error {
	if swag.IsZero(m.VirtualDataPoints) { // not required
		return nil
	}

	for i := 0; i < len(m.VirtualDataPoints); i++ {
		if swag.IsZero(m.VirtualDataPoints[i]) { // not required
			continue
		}

		if m.VirtualDataPoints[i] != nil {
			if err := m.VirtualDataPoints[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("virtualDataPoints" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this pie chart info based on the context it is used
func (m *PieChartInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCounters(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDataPoints(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePieChartItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateVirtualDataPoints(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PieChartInfo) contextValidateCounters(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Counters); i++ {

		if m.Counters[i] != nil {
			if err := m.Counters[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("counters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) contextValidateDataPoints(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.DataPoints); i++ {

		if m.DataPoints[i] != nil {
			if err := m.DataPoints[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("dataPoints" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) contextValidatePieChartItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.PieChartItems); i++ {

		if m.PieChartItems[i] != nil {
			if err := m.PieChartItems[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("pieChartItems" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *PieChartInfo) contextValidateVirtualDataPoints(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.VirtualDataPoints); i++ {

		if m.VirtualDataPoints[i] != nil {
			if err := m.VirtualDataPoints[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("virtualDataPoints" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *PieChartInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PieChartInfo) UnmarshalBinary(b []byte) error {
	var res PieChartInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
