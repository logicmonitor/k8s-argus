// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RestEventSourceFilter rest event source filter
//
// swagger:model RestEventSourceFilter
type RestEventSourceFilter struct {

	// comment
	Comment string `json:"comment,omitempty"`

	// id
	ID int32 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// operator
	Operator string `json:"operator,omitempty"`

	// value
	Value string `json:"value,omitempty"`
}

// Validate validates this rest event source filter
func (m *RestEventSourceFilter) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this rest event source filter based on context it is used
func (m *RestEventSourceFilter) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RestEventSourceFilter) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RestEventSourceFilter) UnmarshalBinary(b []byte) error {
	var res RestEventSourceFilter
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
