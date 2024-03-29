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

// InheritanceProp inheritance prop
//
// swagger:model InheritanceProp
type InheritanceProp struct {

	// fullpath
	// Read Only: true
	Fullpath string `json:"fullpath,omitempty"`

	// id
	// Read Only: true
	ID int32 `json:"id,omitempty"`

	// type
	// Read Only: true
	Type string `json:"type,omitempty"`

	// value
	// Read Only: true
	Value string `json:"value,omitempty"`
}

// Validate validates this inheritance prop
func (m *InheritanceProp) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this inheritance prop based on the context it is used
func (m *InheritanceProp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFullpath(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateType(ctx, formats); err != nil {
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

func (m *InheritanceProp) contextValidateFullpath(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "fullpath", "body", string(m.Fullpath)); err != nil {
		return err
	}

	return nil
}

func (m *InheritanceProp) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", int32(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *InheritanceProp) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "type", "body", string(m.Type)); err != nil {
		return err
	}

	return nil
}

func (m *InheritanceProp) contextValidateValue(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "value", "body", string(m.Value)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InheritanceProp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InheritanceProp) UnmarshalBinary(b []byte) error {
	var res InheritanceProp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
