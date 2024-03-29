// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AlertWidgetData alert widget data
//
// swagger:model AlertWidgetData
type AlertWidgetData struct {
	titleField string

	// items
	// Read Only: true
	Items []*Alert `json:"items,omitempty"`

	// search Id
	// Read Only: true
	SearchID string `json:"searchId,omitempty"`

	// total
	// Read Only: true
	Total int32 `json:"total,omitempty"`
}

// Title gets the title of this subtype
func (m *AlertWidgetData) Title() string {
	return m.titleField
}

// SetTitle sets the title of this subtype
func (m *AlertWidgetData) SetTitle(val string) {
	m.titleField = val
}

// Type gets the type of this subtype
func (m *AlertWidgetData) Type() string {
	return "alert"
}

// SetType sets the type of this subtype
func (m *AlertWidgetData) SetType(val string) {
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *AlertWidgetData) UnmarshalJSON(raw []byte) error {
	var data struct {

		// items
		// Read Only: true
		Items []*Alert `json:"items,omitempty"`

		// search Id
		// Read Only: true
		SearchID string `json:"searchId,omitempty"`

		// total
		// Read Only: true
		Total int32 `json:"total,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		Title string `json:"title,omitempty"`

		Type string `json:"type,omitempty"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result AlertWidgetData

	result.titleField = base.Title

	if base.Type != result.Type() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid type value: %q", base.Type)
	}

	result.Items = data.Items
	result.SearchID = data.SearchID
	result.Total = data.Total

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m AlertWidgetData) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// items
		// Read Only: true
		Items []*Alert `json:"items,omitempty"`

		// search Id
		// Read Only: true
		SearchID string `json:"searchId,omitempty"`

		// total
		// Read Only: true
		Total int32 `json:"total,omitempty"`
	}{

		Items: m.Items,

		SearchID: m.SearchID,

		Total: m.Total,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		Title string `json:"title,omitempty"`

		Type string `json:"type,omitempty"`
	}{

		Title: m.Title(),

		Type: m.Type(),
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this alert widget data
func (m *AlertWidgetData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertWidgetData) validateItems(formats strfmt.Registry) error {

	if swag.IsZero(m.Items) { // not required
		return nil
	}

	for i := 0; i < len(m.Items); i++ {
		if swag.IsZero(m.Items[i]) { // not required
			continue
		}

		if m.Items[i] != nil {
			if err := m.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this alert widget data based on the context it is used
func (m *AlertWidgetData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSearchID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTotal(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertWidgetData) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "type", "body", string(m.Type())); err != nil {
		return err
	}

	return nil
}

func (m *AlertWidgetData) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "items", "body", []*Alert(m.Items)); err != nil {
		return err
	}

	for i := 0; i < len(m.Items); i++ {

		if m.Items[i] != nil {
			if err := m.Items[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *AlertWidgetData) contextValidateSearchID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "searchId", "body", string(m.SearchID)); err != nil {
		return err
	}

	return nil
}

func (m *AlertWidgetData) contextValidateTotal(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "total", "body", int32(m.Total)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertWidgetData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertWidgetData) UnmarshalBinary(b []byte) error {
	var res AlertWidgetData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
