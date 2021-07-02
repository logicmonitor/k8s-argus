// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ReportGroup report group
// swagger:model ReportGroup
type ReportGroup struct {

	// description
	Description string `json:"description,omitempty"`

	// id
	ID int32 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// reports count
	ReportsCount int32 `json:"reportsCount,omitempty"`

	// user permission
	UserPermission string `json:"userPermission,omitempty"`
}

// Validate validates this report group
func (m *ReportGroup) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ReportGroup) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportGroup) UnmarshalBinary(b []byte) error {
	var res ReportGroup
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
