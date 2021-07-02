// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Authentication authentication
// swagger:discriminator Authentication type
type Authentication interface {
	runtime.Validatable

	// password
	// Required: true
	Password() *string
	SetPassword(*string)

	// type
	// Required: true
	Type() string
	SetType(string)

	// user name
	// Required: true
	UserName() *string
	SetUserName(*string)
}

type authentication struct {
	passwordField *string

	typeField string

	userNameField *string
}

// Password gets the password of this polymorphic type
func (m *authentication) Password() *string {
	return m.passwordField
}

// SetPassword sets the password of this polymorphic type
func (m *authentication) SetPassword(val *string) {
	m.passwordField = val
}

// Type gets the type of this polymorphic type
func (m *authentication) Type() string {
	return "Authentication"
}

// SetType sets the type of this polymorphic type
func (m *authentication) SetType(val string) {
}

// UserName gets the user name of this polymorphic type
func (m *authentication) UserName() *string {
	return m.userNameField
}

// SetUserName sets the user name of this polymorphic type
func (m *authentication) SetUserName(val *string) {
	m.userNameField = val
}

// UnmarshalAuthenticationSlice unmarshals polymorphic slices of Authentication
func UnmarshalAuthenticationSlice(reader io.Reader, consumer runtime.Consumer) ([]Authentication, error) {
	var elements []json.RawMessage
	if err := consumer.Consume(reader, &elements); err != nil {
		return nil, err
	}

	var result []Authentication
	for _, element := range elements {
		obj, err := unmarshalAuthentication(element, consumer)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

// UnmarshalAuthentication unmarshals polymorphic Authentication
func UnmarshalAuthentication(reader io.Reader, consumer runtime.Consumer) (Authentication, error) {
	// we need to read this twice, so first into a buffer
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return unmarshalAuthentication(data, consumer)
}

func unmarshalAuthentication(data []byte, consumer runtime.Consumer) (Authentication, error) {
	buf := bytes.NewBuffer(data)
	buf2 := bytes.NewBuffer(data)

	// the first time this is read is to fetch the value of the type property.
	var getType struct {
		Type string `json:"type"`
	}
	if err := consumer.Consume(buf, &getType); err != nil {
		return nil, err
	}

	if err := validate.RequiredString("type", "body", getType.Type); err != nil {
		return nil, err
	}

	// The value of type is used to determine which type to create and unmarshal the data into
	switch getType.Type {
	case "Authentication":
		var result authentication
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case "basic":
		var result BasicAuthentication
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	case "ntlm":
		var result NTLMAuthentication
		if err := consumer.Consume(buf2, &result); err != nil {
			return nil, err
		}
		return &result, nil

	}
	return nil, errors.New(422, "invalid type value: %q", getType.Type)
}

// Validate validates this authentication
func (m *authentication) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *authentication) validatePassword(formats strfmt.Registry) error {
	if err := validate.Required("password", "body", m.Password()); err != nil {
		return err
	}

	return nil
}

func (m *authentication) validateUserName(formats strfmt.Registry) error {
	if err := validate.Required("userName", "body", m.UserName()); err != nil {
		return err
	}

	return nil
}
