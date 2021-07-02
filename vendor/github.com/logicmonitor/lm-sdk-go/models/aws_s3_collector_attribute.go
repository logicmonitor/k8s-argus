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

// AwsS3CollectorAttribute aws s3 collector attribute
// swagger:model AwsS3CollectorAttribute
type AwsS3CollectorAttribute struct {
	AwsS3CollectorAttributeAllOf1
}

// Name gets the name of this subtype
func (m *AwsS3CollectorAttribute) Name() string {
	return "awss3"
}

// SetName sets the name of this subtype
func (m *AwsS3CollectorAttribute) SetName(val string) {
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *AwsS3CollectorAttribute) UnmarshalJSON(raw []byte) error {
	var data struct {
		AwsS3CollectorAttributeAllOf1
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

	var result AwsS3CollectorAttribute

	if base.Name != result.Name() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid name value: %q", base.Name)
	}

	result.AwsS3CollectorAttributeAllOf1 = data.AwsS3CollectorAttributeAllOf1

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m AwsS3CollectorAttribute) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {
		AwsS3CollectorAttributeAllOf1
	}{

		AwsS3CollectorAttributeAllOf1: m.AwsS3CollectorAttributeAllOf1,
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

// Validate validates this aws s3 collector attribute
func (m *AwsS3CollectorAttribute) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with AwsS3CollectorAttributeAllOf1

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *AwsS3CollectorAttribute) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AwsS3CollectorAttribute) UnmarshalBinary(b []byte) error {
	var res AwsS3CollectorAttribute
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// AwsS3CollectorAttributeAllOf1 aws s3 collector attribute all of1
// swagger:model AwsS3CollectorAttributeAllOf1
type AwsS3CollectorAttributeAllOf1 interface{}
