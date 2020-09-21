// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// MapItemInfo map item info
// swagger:model MapItemInfo
type MapItemInfo struct {

	// active status
	// Read Only: true
	ActiveStatus string `json:"activeStatus,omitempty"`

	// alert status
	// Read Only: true
	AlertStatus string `json:"alertStatus,omitempty"`

	// description
	// Read Only: true
	Description string `json:"description,omitempty"`

	// display name
	// Read Only: true
	DisplayName string `json:"displayName,omitempty"`

	// formatted location
	// Read Only: true
	FormattedLocation string `json:"formattedLocation,omitempty"`

	// id
	// Read Only: true
	ID int32 `json:"id,omitempty"`

	// latitude
	// Read Only: true
	Latitude string `json:"latitude,omitempty"`

	// location
	// Read Only: true
	Location string `json:"location,omitempty"`

	// longitude
	// Read Only: true
	Longitude string `json:"longitude,omitempty"`

	// name
	// Read Only: true
	Name string `json:"name,omitempty"`

	// sdt status
	// Read Only: true
	SDTStatus string `json:"sdtStatus,omitempty"`

	// sub type
	// Read Only: true
	SubType string `json:"subType,omitempty"`

	// type
	// Read Only: true
	Type string `json:"type,omitempty"`
}

// Validate validates this map item info
func (m *MapItemInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MapItemInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MapItemInfo) UnmarshalBinary(b []byte) error {
	var res MapItemInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}