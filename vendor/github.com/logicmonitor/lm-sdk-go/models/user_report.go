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

// UserReport user report
//
// swagger:model UserReport
type UserReport struct {
	customReportTypeIdField int32

	customReportTypeNameField string

	deliveryField string

	descriptionField string

	enableViewAsOtherUserField *bool

	formatField string

	groupIdField int32

	idField int32

	lastGenerateOnField int64

	lastGeneratePagesField int32

	lastGenerateSizeField int64

	lastmodifyUserIdField int32

	lastmodifyUserNameField string

	nameField *string

	recipientsField []*ReportRecipient

	reportLinkNumField int32

	scheduleField string

	scheduleTimezoneField string

	userPermissionField string

	// The columns displayed in the report
	Columns []*DynamicColumn `json:"columns,omitempty"`

	// The sort by method
	SortedBy string `json:"sortedBy,omitempty"`

	// The filter for the report
	UserFilter *UserFilter `json:"userFilter,omitempty"`
}

// CustomReportTypeID gets the custom report type Id of this subtype
func (m *UserReport) CustomReportTypeID() int32 {
	return m.customReportTypeIdField
}

// SetCustomReportTypeID sets the custom report type Id of this subtype
func (m *UserReport) SetCustomReportTypeID(val int32) {
	m.customReportTypeIdField = val
}

// CustomReportTypeName gets the custom report type name of this subtype
func (m *UserReport) CustomReportTypeName() string {
	return m.customReportTypeNameField
}

// SetCustomReportTypeName sets the custom report type name of this subtype
func (m *UserReport) SetCustomReportTypeName(val string) {
	m.customReportTypeNameField = val
}

// Delivery gets the delivery of this subtype
func (m *UserReport) Delivery() string {
	return m.deliveryField
}

// SetDelivery sets the delivery of this subtype
func (m *UserReport) SetDelivery(val string) {
	m.deliveryField = val
}

// Description gets the description of this subtype
func (m *UserReport) Description() string {
	return m.descriptionField
}

// SetDescription sets the description of this subtype
func (m *UserReport) SetDescription(val string) {
	m.descriptionField = val
}

// EnableViewAsOtherUser gets the enable view as other user of this subtype
func (m *UserReport) EnableViewAsOtherUser() *bool {
	return m.enableViewAsOtherUserField
}

// SetEnableViewAsOtherUser sets the enable view as other user of this subtype
func (m *UserReport) SetEnableViewAsOtherUser(val *bool) {
	m.enableViewAsOtherUserField = val
}

// Format gets the format of this subtype
func (m *UserReport) Format() string {
	return m.formatField
}

// SetFormat sets the format of this subtype
func (m *UserReport) SetFormat(val string) {
	m.formatField = val
}

// GroupID gets the group Id of this subtype
func (m *UserReport) GroupID() int32 {
	return m.groupIdField
}

// SetGroupID sets the group Id of this subtype
func (m *UserReport) SetGroupID(val int32) {
	m.groupIdField = val
}

// ID gets the id of this subtype
func (m *UserReport) ID() int32 {
	return m.idField
}

// SetID sets the id of this subtype
func (m *UserReport) SetID(val int32) {
	m.idField = val
}

// LastGenerateOn gets the last generate on of this subtype
func (m *UserReport) LastGenerateOn() int64 {
	return m.lastGenerateOnField
}

// SetLastGenerateOn sets the last generate on of this subtype
func (m *UserReport) SetLastGenerateOn(val int64) {
	m.lastGenerateOnField = val
}

// LastGeneratePages gets the last generate pages of this subtype
func (m *UserReport) LastGeneratePages() int32 {
	return m.lastGeneratePagesField
}

// SetLastGeneratePages sets the last generate pages of this subtype
func (m *UserReport) SetLastGeneratePages(val int32) {
	m.lastGeneratePagesField = val
}

// LastGenerateSize gets the last generate size of this subtype
func (m *UserReport) LastGenerateSize() int64 {
	return m.lastGenerateSizeField
}

// SetLastGenerateSize sets the last generate size of this subtype
func (m *UserReport) SetLastGenerateSize(val int64) {
	m.lastGenerateSizeField = val
}

// LastmodifyUserID gets the lastmodify user Id of this subtype
func (m *UserReport) LastmodifyUserID() int32 {
	return m.lastmodifyUserIdField
}

// SetLastmodifyUserID sets the lastmodify user Id of this subtype
func (m *UserReport) SetLastmodifyUserID(val int32) {
	m.lastmodifyUserIdField = val
}

// LastmodifyUserName gets the lastmodify user name of this subtype
func (m *UserReport) LastmodifyUserName() string {
	return m.lastmodifyUserNameField
}

// SetLastmodifyUserName sets the lastmodify user name of this subtype
func (m *UserReport) SetLastmodifyUserName(val string) {
	m.lastmodifyUserNameField = val
}

// Name gets the name of this subtype
func (m *UserReport) Name() *string {
	return m.nameField
}

// SetName sets the name of this subtype
func (m *UserReport) SetName(val *string) {
	m.nameField = val
}

// Recipients gets the recipients of this subtype
func (m *UserReport) Recipients() []*ReportRecipient {
	return m.recipientsField
}

// SetRecipients sets the recipients of this subtype
func (m *UserReport) SetRecipients(val []*ReportRecipient) {
	m.recipientsField = val
}

// ReportLinkNum gets the report link num of this subtype
func (m *UserReport) ReportLinkNum() int32 {
	return m.reportLinkNumField
}

// SetReportLinkNum sets the report link num of this subtype
func (m *UserReport) SetReportLinkNum(val int32) {
	m.reportLinkNumField = val
}

// Schedule gets the schedule of this subtype
func (m *UserReport) Schedule() string {
	return m.scheduleField
}

// SetSchedule sets the schedule of this subtype
func (m *UserReport) SetSchedule(val string) {
	m.scheduleField = val
}

// ScheduleTimezone gets the schedule timezone of this subtype
func (m *UserReport) ScheduleTimezone() string {
	return m.scheduleTimezoneField
}

// SetScheduleTimezone sets the schedule timezone of this subtype
func (m *UserReport) SetScheduleTimezone(val string) {
	m.scheduleTimezoneField = val
}

// Type gets the type of this subtype
func (m *UserReport) Type() string {
	return "User"
}

// SetType sets the type of this subtype
func (m *UserReport) SetType(val string) {
}

// UserPermission gets the user permission of this subtype
func (m *UserReport) UserPermission() string {
	return m.userPermissionField
}

// SetUserPermission sets the user permission of this subtype
func (m *UserReport) SetUserPermission(val string) {
	m.userPermissionField = val
}

// UnmarshalJSON unmarshals this object with a polymorphic type from a JSON structure
func (m *UserReport) UnmarshalJSON(raw []byte) error {
	var data struct {

		// The columns displayed in the report
		Columns []*DynamicColumn `json:"columns,omitempty"`

		// The sort by method
		SortedBy string `json:"sortedBy,omitempty"`

		// The filter for the report
		UserFilter *UserFilter `json:"userFilter,omitempty"`
	}
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&data); err != nil {
		return err
	}

	var base struct {
		/* Just the base type fields. Used for unmashalling polymorphic types.*/

		CustomReportTypeID int32 `json:"customReportTypeId,omitempty"`

		CustomReportTypeName string `json:"customReportTypeName,omitempty"`

		Delivery string `json:"delivery,omitempty"`

		Description string `json:"description,omitempty"`

		EnableViewAsOtherUser *bool `json:"enableViewAsOtherUser,omitempty"`

		Format string `json:"format,omitempty"`

		GroupID int32 `json:"groupId,omitempty"`

		ID int32 `json:"id,omitempty"`

		LastGenerateOn int64 `json:"lastGenerateOn,omitempty"`

		LastGeneratePages int32 `json:"lastGeneratePages,omitempty"`

		LastGenerateSize int64 `json:"lastGenerateSize,omitempty"`

		LastmodifyUserID int32 `json:"lastmodifyUserId,omitempty"`

		LastmodifyUserName string `json:"lastmodifyUserName,omitempty"`

		Name *string `json:"name"`

		Recipients []*ReportRecipient `json:"recipients,omitempty"`

		ReportLinkNum int32 `json:"reportLinkNum,omitempty"`

		Schedule string `json:"schedule,omitempty"`

		ScheduleTimezone string `json:"scheduleTimezone,omitempty"`

		Type string `json:"type"`

		UserPermission string `json:"userPermission,omitempty"`
	}
	buf = bytes.NewBuffer(raw)
	dec = json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&base); err != nil {
		return err
	}

	var result UserReport

	result.customReportTypeIdField = base.CustomReportTypeID

	result.customReportTypeNameField = base.CustomReportTypeName

	result.deliveryField = base.Delivery

	result.descriptionField = base.Description

	result.enableViewAsOtherUserField = base.EnableViewAsOtherUser

	result.formatField = base.Format

	result.groupIdField = base.GroupID

	result.idField = base.ID

	result.lastGenerateOnField = base.LastGenerateOn

	result.lastGeneratePagesField = base.LastGeneratePages

	result.lastGenerateSizeField = base.LastGenerateSize

	result.lastmodifyUserIdField = base.LastmodifyUserID

	result.lastmodifyUserNameField = base.LastmodifyUserName

	result.nameField = base.Name

	result.recipientsField = base.Recipients

	result.reportLinkNumField = base.ReportLinkNum

	result.scheduleField = base.Schedule

	result.scheduleTimezoneField = base.ScheduleTimezone

	if base.Type != result.Type() {
		/* Not the type we're looking for. */
		return errors.New(422, "invalid type value: %q", base.Type)
	}
	result.userPermissionField = base.UserPermission

	result.Columns = data.Columns
	result.SortedBy = data.SortedBy
	result.UserFilter = data.UserFilter

	*m = result

	return nil
}

// MarshalJSON marshals this object with a polymorphic type to a JSON structure
func (m UserReport) MarshalJSON() ([]byte, error) {
	var b1, b2, b3 []byte
	var err error
	b1, err = json.Marshal(struct {

		// The columns displayed in the report
		Columns []*DynamicColumn `json:"columns,omitempty"`

		// The sort by method
		SortedBy string `json:"sortedBy,omitempty"`

		// The filter for the report
		UserFilter *UserFilter `json:"userFilter,omitempty"`
	}{

		Columns: m.Columns,

		SortedBy: m.SortedBy,

		UserFilter: m.UserFilter,
	})
	if err != nil {
		return nil, err
	}
	b2, err = json.Marshal(struct {
		CustomReportTypeID int32 `json:"customReportTypeId,omitempty"`

		CustomReportTypeName string `json:"customReportTypeName,omitempty"`

		Delivery string `json:"delivery,omitempty"`

		Description string `json:"description,omitempty"`

		EnableViewAsOtherUser *bool `json:"enableViewAsOtherUser,omitempty"`

		Format string `json:"format,omitempty"`

		GroupID int32 `json:"groupId,omitempty"`

		ID int32 `json:"id,omitempty"`

		LastGenerateOn int64 `json:"lastGenerateOn,omitempty"`

		LastGeneratePages int32 `json:"lastGeneratePages,omitempty"`

		LastGenerateSize int64 `json:"lastGenerateSize,omitempty"`

		LastmodifyUserID int32 `json:"lastmodifyUserId,omitempty"`

		LastmodifyUserName string `json:"lastmodifyUserName,omitempty"`

		Name *string `json:"name"`

		Recipients []*ReportRecipient `json:"recipients,omitempty"`

		ReportLinkNum int32 `json:"reportLinkNum,omitempty"`

		Schedule string `json:"schedule,omitempty"`

		ScheduleTimezone string `json:"scheduleTimezone,omitempty"`

		Type string `json:"type"`

		UserPermission string `json:"userPermission,omitempty"`
	}{

		CustomReportTypeID: m.CustomReportTypeID(),

		CustomReportTypeName: m.CustomReportTypeName(),

		Delivery: m.Delivery(),

		Description: m.Description(),

		EnableViewAsOtherUser: m.EnableViewAsOtherUser(),

		Format: m.Format(),

		GroupID: m.GroupID(),

		ID: m.ID(),

		LastGenerateOn: m.LastGenerateOn(),

		LastGeneratePages: m.LastGeneratePages(),

		LastGenerateSize: m.LastGenerateSize(),

		LastmodifyUserID: m.LastmodifyUserID(),

		LastmodifyUserName: m.LastmodifyUserName(),

		Name: m.Name(),

		Recipients: m.Recipients(),

		ReportLinkNum: m.ReportLinkNum(),

		Schedule: m.Schedule(),

		ScheduleTimezone: m.ScheduleTimezone(),

		Type: m.Type(),

		UserPermission: m.UserPermission(),
	})
	if err != nil {
		return nil, err
	}

	return swag.ConcatJSON(b1, b2, b3), nil
}

// Validate validates this user report
func (m *UserReport) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRecipients(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateColumns(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserFilter(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UserReport) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name()); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) validateRecipients(formats strfmt.Registry) error {

	if swag.IsZero(m.Recipients()) { // not required
		return nil
	}

	for i := 0; i < len(m.Recipients()); i++ {
		if swag.IsZero(m.recipientsField[i]) { // not required
			continue
		}

		if m.recipientsField[i] != nil {
			if err := m.recipientsField[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("recipients" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *UserReport) validateColumns(formats strfmt.Registry) error {

	if swag.IsZero(m.Columns) { // not required
		return nil
	}

	for i := 0; i < len(m.Columns); i++ {
		if swag.IsZero(m.Columns[i]) { // not required
			continue
		}

		if m.Columns[i] != nil {
			if err := m.Columns[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("columns" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *UserReport) validateUserFilter(formats strfmt.Registry) error {

	if swag.IsZero(m.UserFilter) { // not required
		return nil
	}

	if m.UserFilter != nil {
		if err := m.UserFilter.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("userFilter")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this user report based on the context it is used
func (m *UserReport) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCustomReportTypeID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCustomReportTypeName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateEnableViewAsOtherUser(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastGenerateOn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastGeneratePages(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastGenerateSize(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastmodifyUserID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLastmodifyUserName(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRecipients(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateReportLinkNum(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUserPermission(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateColumns(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUserFilter(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UserReport) contextValidateCustomReportTypeID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "customReportTypeId", "body", int32(m.CustomReportTypeID())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateCustomReportTypeName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "customReportTypeName", "body", string(m.CustomReportTypeName())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateEnableViewAsOtherUser(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "enableViewAsOtherUser", "body", m.EnableViewAsOtherUser()); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "id", "body", int32(m.ID())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateLastGenerateOn(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "lastGenerateOn", "body", int64(m.LastGenerateOn())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateLastGeneratePages(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "lastGeneratePages", "body", int32(m.LastGeneratePages())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateLastGenerateSize(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "lastGenerateSize", "body", int64(m.LastGenerateSize())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateLastmodifyUserID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "lastmodifyUserId", "body", int32(m.LastmodifyUserID())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateLastmodifyUserName(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "lastmodifyUserName", "body", string(m.LastmodifyUserName())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateRecipients(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Recipients()); i++ {

		if m.recipientsField[i] != nil {
			if err := m.recipientsField[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("recipients" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *UserReport) contextValidateReportLinkNum(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "reportLinkNum", "body", int32(m.ReportLinkNum())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateUserPermission(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "userPermission", "body", string(m.UserPermission())); err != nil {
		return err
	}

	return nil
}

func (m *UserReport) contextValidateColumns(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Columns); i++ {

		if m.Columns[i] != nil {
			if err := m.Columns[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("columns" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *UserReport) contextValidateUserFilter(ctx context.Context, formats strfmt.Registry) error {

	if m.UserFilter != nil {
		if err := m.UserFilter.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("userFilter")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *UserReport) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserReport) UnmarshalBinary(b []byte) error {
	var res UserReport
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
