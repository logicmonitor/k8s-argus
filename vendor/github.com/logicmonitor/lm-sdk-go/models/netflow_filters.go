// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NetflowFilters netflow filters
// swagger:model NetflowFilters
type NetflowFilters struct {

	// direction
	Direction string `json:"direction,omitempty"`

	// if idx
	IfIdx int32 `json:"ifIdx,omitempty"`

	// if name
	IfName string `json:"ifName,omitempty"`

	// node a
	NodeA string `json:"nodeA,omitempty"`

	// node b
	NodeB string `json:"nodeB,omitempty"`

	// ports
	Ports string `json:"ports,omitempty"`

	// protocol
	Protocol string `json:"protocol,omitempty"`

	// qos type
	QosType string `json:"qosType,omitempty"`

	// top
	Top int32 `json:"top,omitempty"`
}

// Validate validates this netflow filters
func (m *NetflowFilters) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NetflowFilters) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetflowFilters) UnmarshalBinary(b []byte) error {
	var res NetflowFilters
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
