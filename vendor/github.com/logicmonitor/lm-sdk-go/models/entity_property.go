// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// EntityProperty entity property
// swagger:model EntityProperty
type EntityProperty struct {

	// inherit list
	InheritList []*InheritanceProp `json:"inheritList,omitempty"`

	// name
	// Read Only: true
	Name string `json:"name,omitempty"`

	// type
	// Read Only: true
	Type string `json:"type,omitempty"`

	// value
	// Read Only: true
	Value string `json:"value,omitempty"`
}

// Validate validates this entity property
func (m *EntityProperty) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInheritList(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EntityProperty) validateInheritList(formats strfmt.Registry) error {
	if swag.IsZero(m.InheritList) { // not required
		return nil
	}

	for i := 0; i < len(m.InheritList); i++ {
		if swag.IsZero(m.InheritList[i]) { // not required
			continue
		}

		if m.InheritList[i] != nil {
			if err := m.InheritList[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("inheritList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *EntityProperty) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EntityProperty) UnmarshalBinary(b []byte) error {
	var res EntityProperty
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
