package device

import (
	"context"
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
)

// Manager implements types.DeviceManager
type Manager struct {
	*types.Base
	*builder.Builder
	ControllerClient api.CollectorSetControllerClient
}

func buildDevice(c *config.Config, client api.CollectorSetControllerClient, options ...types.DeviceOption) *models.Device {
	hostGroupIds := "1"
	propertyName := constants.K8sClusterNamePropertyKey
	// use the copy value
	clusterName := c.ClusterName
	device := &models.Device{
		CustomProperties: []*models.NameAndValue{
			{
				Name:  &propertyName,
				Value: &clusterName,
			},
		},
		DisableAlerting: c.DisableAlerting,
		HostGroupIds:    &hostGroupIds,
		DeviceType:      constants.K8sDeviceType,
	}

	for _, option := range options {
		option(device)
	}

	reply, err := client.CollectorID(context.Background(), &api.CollectorIDRequest{})
	if err != nil {
		log.Errorf("Failed to get collector ID: %v", err)
	} else {
		log.Infof("Using collector ID %d for %q", reply.Id, *device.DisplayName)
		device.PreferredCollectorID = &reply.Id
	}

	return device
}

// renameAndAddDevice rename display name and then add the device
func (m *Manager) renameAndAddDevice(device *models.Device) (*models.Device, error) {
	resourceName := m.GetPropertyValue(device, constants.K8sResourceNamePropertyKey)
	if resourceName == "" {
		resourceName = *device.DisplayName
	}
	if resourceName == "" {
		return nil, fmt.Errorf("get device(%s) resource name failed", *device.DisplayName)
	}
	renameResourceName := fmt.Sprintf("%s-%s", resourceName, m.Config().ClusterName)
	existDevice, err := m.FindByDisplayName(renameResourceName)
	if err != nil {
		log.Warnf("Get device(%s) failed, err: %v", resourceName, err)
	}
	if existDevice != nil {
		if m.Config().ClusterName == m.GetPropertyValue(existDevice, constants.K8sClusterNamePropertyKey) {
			device.DisplayName = existDevice.DisplayName
			return m.updateAndReplace(existDevice.ID, device)
		}
		return nil, fmt.Errorf("exist displayName: %s", renameResourceName)
	}
	log.Infof("Rename device: %s -> %s", *device.DisplayName, renameResourceName)
	device.DisplayName = &renameResourceName
	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(device)
	restResponse, err := m.LMClient.LM.AddDevice(params)
	if err != nil {
		return nil, err
	}
	return restResponse.Payload, nil
}

// GetPropertyValue get device property value by property name
func (m *Manager) GetPropertyValue(device *models.Device, propertyName string) string {
	if device == nil {
		return ""
	}
	if len(device.CustomProperties) > 0 {
		for _, cp := range device.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}
	if len(device.SystemProperties) > 0 {
		for _, cp := range device.SystemProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}
	return ""
}

func (m *Manager) updateAndReplace(id int32, device *models.Device) (*models.Device, error) {
	opType := "replace"
	params := lm.NewUpdateDeviceParams()
	params.SetID(id)
	params.SetBody(device)
	params.SetOpType(&opType)

	restResponse, err := m.LMClient.LM.UpdateDevice(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)

	return restResponse.Payload, nil
}

// FindByDisplayName implements types.DeviceManager.
func (m *Manager) FindByDisplayName(name string) (*models.Device, error) {
	filter := fmt.Sprintf("displayName:\"%s\"", name)
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)
	if restResponse.Payload.Total == 1 {
		return restResponse.Payload.Items[0], nil
	}

	return nil, nil
}

// FindByDisplayNames implements types.DeviceManager.
func (m *Manager) FindByDisplayNames(displayNames ...string) ([]*models.Device, error) {
	if len(displayNames) == 0 {
		return []*models.Device{}, nil
	}
	filter := fmt.Sprintf("displayName:\"%s\"", strings.Join(displayNames, "\"|\""))
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)
	return restResponse.Payload.Items, nil
}

// FindByDisplayNameAndClusterName implements types.DeviceManager.
func (m *Manager) FindByDisplayNameAndClusterName(displayName string) (*models.Device, error) {
	displayNameWithClusterName := fmt.Sprintf("%s-%s", displayName, m.Config().ClusterName)
	devices, err := m.FindByDisplayNames(displayName, displayNameWithClusterName)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if m.Config().ClusterName == m.GetPropertyValue(device, constants.K8sClusterNamePropertyKey) {
			return device, nil
		}
	}
	return nil, nil
}

// Add implements types.DeviceManager.
func (m *Manager) Add(options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), m.ControllerClient, options...)
	log.Debugf("%#v", device)

	existDevice, err := m.checkAndUpdateExistingDevice(device)
	if err != nil {
		return nil, err
	}
	if existDevice != nil {
		return existDevice, nil
	}

	// add the new device
	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(device)
	restResponse, err := m.LMClient.LM.AddDevice(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)
	return restResponse.Payload, nil
}

func (m *Manager) checkAndUpdateExistingDevice(device *models.Device) (*models.Device, error) {
	clusterName := m.Config().ClusterName
	displayNameWithClusterName := fmt.Sprintf("%s-%s", *device.DisplayName, clusterName)
	devices, err := m.FindByDisplayNames(*device.DisplayName, displayNameWithClusterName)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", devices)
	if len(devices) == 0 {
		return nil, nil
	}
	var existDevice *models.Device
	for _, device := range devices {
		if clusterName == m.GetPropertyValue(device, constants.K8sClusterNamePropertyKey) {
			existDevice = device
			break
		}
	}
	if existDevice != nil {
		// exist the device, just update it
		*device.DisplayName = *existDevice.DisplayName
		updatedDevice, err1 := m.updateAndReplace(existDevice.ID, device)
		if err1 != nil {
			return nil, err1
		}
		log.Infof("Exists the device(%s), just update it", *updatedDevice.DisplayName)
		return updatedDevice, nil
	}

	// exist the duplicate name device, rename displayName and add device
	renamedDevice, err := m.renameAndAddDevice(device)
	if err != nil {
		log.Errorf("rename device failed: %v", err)
		return nil, fmt.Errorf("rename device failed")
	}
	return renamedDevice, nil
}

// UpdateAndReplaceByID implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByID(id int32, options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), m.ControllerClient, options...)
	log.Debugf("%#v", device)

	return m.updateAndReplace(id, device)
}

// UpdateAndReplaceByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByDisplayName(name string, options ...types.DeviceOption) (*models.Device, error) {
	d, err := m.FindByDisplayNameAndClusterName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Warnf("Could not find device %q", name)
		return nil, nil
	}
	options = append(options, m.DisplayName(*d.DisplayName))
	// Update the device.
	device, err := m.UpdateAndReplaceByID(d.ID, options...)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// TODO: this method needs to be removed in DEV-50496

// UpdateAndReplaceFieldByID implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceFieldByID(id int32, field string, options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), m.ControllerClient, options...)
	log.Debugf("%#v", device)

	params := lm.NewPatchDeviceParams()
	params.SetID(id)
	params.SetBody(device)
	opType := "replace"
	params.SetOpType(&opType)
	restResponse, err := m.LMClient.LM.PatchDevice(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)

	return restResponse.Payload, nil
}

// TODO: this method needs to be removed in DEV-50496

// UpdateAndReplaceFieldByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceFieldByDisplayName(name string, field string, options ...types.DeviceOption) (*models.Device, error) {
	d, err := m.FindByDisplayNameAndClusterName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Infof("Could not find device %q", name)
		return nil, nil
	}
	options = append(options, m.DisplayName(*d.DisplayName))
	// Update the device.
	device, err := m.UpdateAndReplaceFieldByID(d.ID, field, options...)
	if err != nil {
		return nil, err
	}

	return device, nil
}

// DeleteByID implements types.DeviceManager.
func (m *Manager) DeleteByID(id int32) error {
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(id)
	_, err := m.LMClient.LM.DeleteDeviceByID(params)
	return err
}

// DeleteByDisplayName implements types.DeviceManager.
func (m *Manager) DeleteByDisplayName(name string) error {
	d, err := m.FindByDisplayNameAndClusterName(name)
	if err != nil {
		return err
	}

	// TODO: Should this return an error?
	if d == nil {
		log.Infof("Could not find device %q", name)
		return nil
	}

	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(d.ID)
	_, err = m.LMClient.LM.DeleteDeviceByID(params)
	return err
}

// Config implements types.DeviceManager.
func (m *Manager) Config() *config.Config {
	return m.Base.Config
}

// GetListByGroupID implements getting all the devices belongs to the group directly
func (m *Manager) GetListByGroupID(groupID int32) ([]*models.Device, error) {
	params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
	params.SetID(groupID)
	fields := "id,name,displayName,customProperties"
	params.SetFields(&fields)
	size := int32(-1)
	params.SetSize(&size)
	restResponse, err := m.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)
	return restResponse.Payload.Items, nil
}
