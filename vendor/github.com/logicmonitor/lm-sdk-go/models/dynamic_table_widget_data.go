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

// DynamicTableWidgetData dynamic table widget data
//
// swagger:model DynamicTableWidgetData
type DynamicTableWidgetData struct {
	titleField string

	// column headers
	// Read Only: true
	ColumnHeaders []*ColumnHeader `json:"columnHeaders,omitempty"`

	// rows
	// Read Only: true
	Rows []*RowData `json:"rows,omitempty"`
}

// Title gets the title of this subtype
func (m *DynamicTableWidgetData) Title() string {
	return m.titleField
}

// SetTitle sets the title of this subtype
func (m *DynamicTableWidgetData) SetTitle(val string) {
	m.titleField = val
}

// Type gets the type of this subtype
func (m *DynamicTableWidgetData) Type() string {
	return "dynamicTable"
}

// SetType sets the type of this subtype
func (m *DynamicTableWidgetData) SetType(val string) {
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *DynamicTableWidgetData) UnmarshalJSON(raw []byte) error {
	var data struct {

		// column headers
		// Read Only: true
		ColumnHeaders []*ColumnHeader `json:"columnHeaders,omitempty"`

		// rows
		// Read Only: true
		Rows []*RowData `json:"rows,omitempty"`
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

	var result DynamicTableWidgetData

	result.titleField = base.Title

	if base.Type != result.Type() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid type value: %q", base.Type)
	}

	result.ColumnHeaders = data.ColumnHeaders
	result.Rows = data.Rows

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m DynamicTableWidgetData) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// column headers
		// Read Only: true
		ColumnHeaders []*ColumnHeader `json:"columnHeaders,omitempty"`

		// rows
		// Read Only: true
		Rows []*RowData `json:"rows,omitempty"`
	}{

		ColumnHeaders: m.ColumnHeaders,

		Rows: m.Rows,
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

// Validate validates this dynamic table widget data
func (m *DynamicTableWidgetData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateColumnHeaders(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRows(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DynamicTableWidgetData) validateColumnHeaders(formats strfmt.Registry) error {

	if swag.IsZero(m.ColumnHeaders) { // not required
		return nil
	}

	for i := 0; i < len(m.ColumnHeaders); i++ {
		if swag.IsZero(m.ColumnHeaders[i]) { // not required
			continue
		}

		if m.ColumnHeaders[i] != nil {
			if err := m.ColumnHeaders[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("columnHeaders" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DynamicTableWidgetData) validateRows(formats strfmt.Registry) error {

	if swag.IsZero(m.Rows) { // not required
		return nil
	}

	for i := 0; i < len(m.Rows); i++ {
		if swag.IsZero(m.Rows[i]) { // not required
			continue
		}

		if m.Rows[i] != nil {
			if err := m.Rows[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("rows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this dynamic table widget data based on the context it is used
func (m *DynamicTableWidgetData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateColumnHeaders(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRows(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DynamicTableWidgetData) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "type", "body", string(m.Type())); err != nil {
		return err
	}

	return nil
}

func (m *DynamicTableWidgetData) contextValidateColumnHeaders(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "columnHeaders", "body", []*ColumnHeader(m.ColumnHeaders)); err != nil {
		return err
	}

	for i := 0; i < len(m.ColumnHeaders); i++ {

		if m.ColumnHeaders[i] != nil {
			if err := m.ColumnHeaders[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("columnHeaders" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DynamicTableWidgetData) contextValidateRows(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "rows", "body", []*RowData(m.Rows)); err != nil {
		return err
	}

	for i := 0; i < len(m.Rows); i++ {

		if m.Rows[i] != nil {
			if err := m.Rows[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("rows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *DynamicTableWidgetData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DynamicTableWidgetData) UnmarshalBinary(b []byte) error {
	var res DynamicTableWidgetData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
