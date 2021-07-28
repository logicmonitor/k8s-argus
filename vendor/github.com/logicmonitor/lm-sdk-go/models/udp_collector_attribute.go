// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UDPCollectorAttribute UDP collector attribute
// swagger:model UDPCollectorAttribute
type UDPCollectorAttribute struct {

	// payload
	Payload string `json:"payload,omitempty"`

	// port
	Port string `json:"port,omitempty"`
}

// Name gets the name of this subtype
func (m *UDPCollectorAttribute) Name() string {
	return "udp"
}

// SetName sets the name of this subtype
func (m *UDPCollectorAttribute) SetName(val string) {
}

// Payload gets the payload of this subtype

// Port gets the port of this subtype

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *UDPCollectorAttribute) UnmarshalJSON(raw []byte) error {
	var data struct {

		// payload
		Payload string `json:"payload,omitempty"`

		// port
		Port string `json:"port,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		Name string `json:"name"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result UDPCollectorAttribute

	if base.Name != result.Name() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid name value: %q", base.Name)
	}

	result.Payload = data.Payload

	result.Port = data.Port

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m UDPCollectorAttribute) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// payload
		Payload string `json:"payload,omitempty"`

		// port
		Port string `json:"port,omitempty"`
	}{

		Payload: m.Payload,

		Port: m.Port,
	},
	)
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		Name string `json:"name"`
	}{

		Name: m.Name(),
	},
	)
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this UDP collector attribute
func (m *UDPCollectorAttribute) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *UDPCollectorAttribute) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UDPCollectorAttribute) UnmarshalBinary(b []byte) error {
	var res UDPCollectorAttribute
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
