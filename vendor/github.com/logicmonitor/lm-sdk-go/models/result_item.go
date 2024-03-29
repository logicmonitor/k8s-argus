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

// ResultItem result item
//
// swagger:model ResultItem
type ResultItem struct {

	// bottom label
	// Read Only: true
	BottomLabel string `json:"bottomLabel,omitempty"`

	// color level
	// Read Only: true
	ColorLevel int32 `json:"colorLevel,omitempty"`

	// value
	// Read Only: true
	Value string `json:"value,omitempty"`
}

// Validate validates this result item
func (m *ResultItem) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this result item based on the context it is used
func (m *ResultItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBottomLabel(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateColorLevel(ctx, formats); err != nil {
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

func (m *ResultItem) contextValidateBottomLabel(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "bottomLabel", "body", string(m.BottomLabel)); err != nil {
		return err
	}

	return nil
}

func (m *ResultItem) contextValidateColorLevel(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "colorLevel", "body", int32(m.ColorLevel)); err != nil {
		return err
	}

	return nil
}

func (m *ResultItem) contextValidateValue(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "value", "body", string(m.Value)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ResultItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResultItem) UnmarshalBinary(b []byte) error {
	var res ResultItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
