// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// NetflowBandwidth netflow bandwidth
// swagger:model NetflowBandwidth
type NetflowBandwidth struct {

	// device display name
	DeviceDisplayName string `json:"deviceDisplayName,omitempty"`

	// receive
	Receive float64 `json:"receive,omitempty"`

	// send
	Send float64 `json:"send,omitempty"`

	// usage
	Usage float64 `json:"usage,omitempty"`
}

// DataType gets the data type of this subtype
func (m *NetflowBandwidth) DataType() string {
	return "bandwidth"
}

// SetDataType sets the data type of this subtype
func (m *NetflowBandwidth) SetDataType(val string) {

}

// DeviceDisplayName gets the device display name of this subtype

// Receive gets the receive of this subtype

// Send gets the send of this subtype

// Usage gets the usage of this subtype

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *NetflowBandwidth) UnmarshalJSON(raw []byte) error {
	var data struct {

		// device display name
		DeviceDisplayName string `json:"deviceDisplayName,omitempty"`

		// receive
		Receive float64 `json:"receive,omitempty"`

		// send
		Send float64 `json:"send,omitempty"`

		// usage
		Usage float64 `json:"usage,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		DataType string `json:"dataType,omitempty"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result NetflowBandwidth

	if base.DataType != result.DataType() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid dataType value: %q", base.DataType)
	}

	result.DeviceDisplayName = data.DeviceDisplayName

	result.Receive = data.Receive

	result.Send = data.Send

	result.Usage = data.Usage

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m NetflowBandwidth) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// device display name
		DeviceDisplayName string `json:"deviceDisplayName,omitempty"`

		// receive
		Receive float64 `json:"receive,omitempty"`

		// send
		Send float64 `json:"send,omitempty"`

		// usage
		Usage float64 `json:"usage,omitempty"`
	}{

		DeviceDisplayName: m.DeviceDisplayName,

		Receive: m.Receive,

		Send: m.Send,

		Usage: m.Usage,
	},
	)
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		DataType string `json:"dataType,omitempty"`
	}{

		DataType: m.DataType(),
	},
	)
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this netflow bandwidth
func (m *NetflowBandwidth) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *NetflowBandwidth) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetflowBandwidth) UnmarshalBinary(b []byte) error {
	var res NetflowBandwidth
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}