// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Usage usage
//
// swagger:model Usage
type Usage struct {

	// Number of AWS resources not monitored with a local Collector
	// Read Only: true
	NumOfAWSDevices int32 `json:"numOfAWSDevices,omitempty"`

	// Number of Azure resources not monitored with a local Collector
	// Read Only: true
	NumOfAzureDevices int32 `json:"numOfAzureDevices,omitempty"`

	// Number of AWS resources monitored with a local Collector
	// Read Only: true
	NumOfCombinedAWSDevices int32 `json:"numOfCombinedAWSDevices,omitempty"`

	// Number of Azure resources monitored with a local Collector
	// Read Only: true
	NumOfCombinedAzureDevices int32 `json:"numOfCombinedAzureDevices,omitempty"`

	// Number of GCP resources monitored with a local Collector
	// Read Only: true
	NumOfCombinedGcpDevices int32 `json:"numOfCombinedGcpDevices,omitempty"`

	// Number of devices with active ConfigSources
	// Read Only: true
	NumOfConfigSourceDevices int32 `json:"numOfConfigSourceDevices,omitempty"`

	// Number of GCP resources
	// Read Only: true
	NumOfGcpDevices int32 `json:"numOfGcpDevices,omitempty"`

	// Number of services (created via LM Service Insight)
	// Read Only: true
	NumOfServices int32 `json:"numOfServices,omitempty"`

	// Number of stopped AWS resources
	// Read Only: true
	NumOfStoppedAWSDevices int32 `json:"numOfStoppedAWSDevices,omitempty"`

	// Number of stopped Azure resources
	// Read Only: true
	NumOfStoppedAzureDevices int32 `json:"numOfStoppedAzureDevices,omitempty"`

	// Number of stopped GCP resources not monitored with a local Collector
	// Read Only: true
	NumOfStoppedGcpDevices int32 `json:"numOfStoppedGcpDevices,omitempty"`

	// Number of terminated AWS resources
	// Read Only: true
	NumOfTerminatedAWSDevices int32 `json:"numOfTerminatedAWSDevices,omitempty"`

	// Number of terminated Azure resources
	// Read Only: true
	NumOfTerminatedAzureDevices int32 `json:"numOfTerminatedAzureDevices,omitempty"`

	// Number of terminated GCP resources
	// Read Only: true
	NumOfTerminatedGcpCloudDevices int32 `json:"numOfTerminatedGcpCloudDevices,omitempty"`

	// Number of websites
	// Read Only: true
	NumOfWebsites int32 `json:"numOfWebsites,omitempty"`

	// Sum of numOfStandardDevices, numOfCombinedAWSDevices, numOfCombinedAzureDevices, and numOfCombinedGCPDevices
	// Read Only: true
	NumberOfDevices int32 `json:"numberOfDevices,omitempty"`

	// Number of monitored Kubernetes Nodes, Pods, and Services
	// Read Only: true
	NumberOfKubernetesDevices int32 `json:"numberOfKubernetesDevices,omitempty"`

	// Number of standard devices
	// Read Only: true
	NumberOfStandardDevices int32 `json:"numberOfStandardDevices,omitempty"`
}

// Validate validates this usage
func (m *Usage) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this usage based on the context it is used
func (m *Usage) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateNumOfAWSDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfAzureDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfCombinedAWSDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfCombinedAzureDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfCombinedGcpDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfConfigSourceDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfGcpDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfServices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfStoppedAWSDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfStoppedAzureDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfStoppedGcpDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfTerminatedAWSDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfTerminatedAzureDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfTerminatedGcpCloudDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumOfWebsites(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumberOfDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumberOfKubernetesDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNumberOfStandardDevices(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Usage) contextValidateNumOfAWSDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfAWSDevices", "body", int32(m.NumOfAWSDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfAzureDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfAzureDevices", "body", int32(m.NumOfAzureDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfCombinedAWSDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfCombinedAWSDevices", "body", int32(m.NumOfCombinedAWSDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfCombinedAzureDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfCombinedAzureDevices", "body", int32(m.NumOfCombinedAzureDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfCombinedGcpDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfCombinedGcpDevices", "body", int32(m.NumOfCombinedGcpDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfConfigSourceDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfConfigSourceDevices", "body", int32(m.NumOfConfigSourceDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfGcpDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfGcpDevices", "body", int32(m.NumOfGcpDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfServices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfServices", "body", int32(m.NumOfServices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfStoppedAWSDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfStoppedAWSDevices", "body", int32(m.NumOfStoppedAWSDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfStoppedAzureDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfStoppedAzureDevices", "body", int32(m.NumOfStoppedAzureDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfStoppedGcpDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfStoppedGcpDevices", "body", int32(m.NumOfStoppedGcpDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfTerminatedAWSDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfTerminatedAWSDevices", "body", int32(m.NumOfTerminatedAWSDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfTerminatedAzureDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfTerminatedAzureDevices", "body", int32(m.NumOfTerminatedAzureDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfTerminatedGcpCloudDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfTerminatedGcpCloudDevices", "body", int32(m.NumOfTerminatedGcpCloudDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumOfWebsites(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numOfWebsites", "body", int32(m.NumOfWebsites)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumberOfDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numberOfDevices", "body", int32(m.NumberOfDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumberOfKubernetesDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numberOfKubernetesDevices", "body", int32(m.NumberOfKubernetesDevices)); err != nil {
		return err
	}

	return nil
}

func (m *Usage) contextValidateNumberOfStandardDevices(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "numberOfStandardDevices", "body", int32(m.NumberOfStandardDevices)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Usage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Usage) UnmarshalBinary(b []byte) error {
	var res Usage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
