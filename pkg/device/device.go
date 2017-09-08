package device

import (
	"fmt"
	"log"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/types"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
)

// Manager implements types.DeviceManager
type Manager struct {
	*types.Base
	*builder.Builder
}

func buildDevice(c *config.Config, options ...types.DeviceOption) *lm.RestDevice {
	device := &lm.RestDevice{
		CustomProperties: []lm.NameAndValue{
			{
				Name:  "auto.clustername",
				Value: c.ClusterName,
			},
		},
		DisableAlerting:      c.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: c.PreferredCollector,
	}

	for _, option := range options {
		option(device)
	}

	return device
}

// FindByName implements types.DeviceManager.
func (m *Manager) FindByName(name string) (*lm.RestDevice, error) {
	return find("name", name, m.LMClient)
}

// FindByDisplayName implements types.DeviceManager.
func (m *Manager) FindByDisplayName(name string) (*lm.RestDevice, error) {
	return find("displayName", name, m.LMClient)
}

// Add implements types.DeviceManager.
func (m *Manager) Add(options ...types.DeviceOption) (*lm.RestDevice, error) {
	device := buildDevice(m.Config(), options...)
	restResponse, apiResponse, err := m.LMClient.AddDevice(*device, false)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return nil, err
	}

	return &restResponse.Data, nil
}

// UpdateAndReplaceByID implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByID(id int32, options ...types.DeviceOption) (*lm.RestDevice, error) {
	device := buildDevice(m.Config(), options...)
	restResponse, apiResponse, err := m.LMClient.UpdateDevice(*device, id, "replace")
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return nil, err
	}

	return &restResponse.Data, nil
}

// UpdateAndReplaceByName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByName(name string, options ...types.DeviceOption) (*lm.RestDevice, error) {
	d, err := m.FindByDisplayName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Printf("Could not find device %q", name)
		return nil, nil
	}

	// Update the device.
	device, err := m.UpdateAndReplaceByID(d.Id, options...)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// UpdateAndReplaceFieldByID implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceFieldByID(id int32, field string, options ...types.DeviceOption) (*lm.RestDevice, error) {
	device := buildDevice(m.Config(), options...)
	restResponse, apiResponse, err := m.LMClient.PatchDeviceById(*device, id, "replace", field)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return nil, err
	}

	return &restResponse.Data, nil
}

// UpdateAndReplaceFieldByName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceFieldByName(name string, field string, options ...types.DeviceOption) (*lm.RestDevice, error) {
	d, err := m.FindByDisplayName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Printf("Could not find device %q", name)
		return nil, nil
	}

	// Update the device.
	device, err := m.UpdateAndReplaceFieldByID(d.Id, field, options...)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// DeleteByID implements types.DeviceManager.
func (m *Manager) DeleteByID(id int32) error {
	restResponse, apiResponse, err := m.LMClient.DeleteDevice(id)
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// DeleteByName implements types.DeviceManager.
func (m *Manager) DeleteByName(name string) error {
	d, err := m.FindByDisplayName(name)
	if err != nil {
		return err
	}

	// TODO: Should this return an error?
	if d == nil {
		log.Printf("Could not find device %q", name)
		return nil
	}

	restResponse, apiResponse, err := m.LMClient.DeleteDevice(d.Id)
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// Config implements types.DeviceManager.
func (m *Manager) Config() *config.Config {
	return m.Base.Config
}

func find(field, name string, client *lm.DefaultApi) (*lm.RestDevice, error) {
	filter := fmt.Sprintf("%s:%s", field, name)
	restResponse, apiResponse, err := client.GetDeviceList("", -1, 0, filter)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return nil, _err
	}
	if restResponse.Data.Total == 1 {
		return &restResponse.Data.Items[0], nil
	}

	return nil, nil
}
