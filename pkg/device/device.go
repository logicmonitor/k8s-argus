package device

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	cscutils "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
)

// Manager implements types.DeviceManager
type Manager struct {
	*types.Base
	*builder.Builder
	DC *devicecache.DeviceCache
}

func buildDevice(c *config.Config, d *models.Device, options ...types.DeviceOption) *models.Device {
	if d == nil {
		hostGroupIds := "1"
		propertyName := constants.K8sClusterNamePropertyKey
		// use the copy value
		clusterName := c.ClusterName
		d = &models.Device{
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
			option(d)
		}

		collectorID := cscutils.GetCollectorID()
		log.Infof("Using collector ID %d for %q", collectorID, *d.DisplayName)
		d.PreferredCollectorID = &collectorID
	} else {
		for _, option := range options {
			option(d)
		}
	}

	return d
}

// checkAndUpdateExistingDevice tries to find and update the devices which needs to be changed
func (m *Manager) checkAndUpdateExistingDevice(device *models.Device) (*models.Device, error) {
	displayNameWithClusterName := fmt.Sprintf("%s-%s", *device.DisplayName, m.Config().ClusterName)
	existingDevices, err := m.FindByDisplayNames(*device.DisplayName, displayNameWithClusterName)
	if err != nil {
		return nil, err
	}
	if len(existingDevices) == 0 {
		return nil, fmt.Errorf("cannot find devices with name: %s", *device.DisplayName)
	}
	for _, existingDevice := range existingDevices {
		clusterName := m.GetPropertyValue(existingDevice, constants.K8sClusterNamePropertyKey)
		if clusterName == m.Config().ClusterName {
			// the device which is not changed will be ignored
			if *existingDevice.Name == *device.Name {
				log.Infof("No changes to device (%s). Ignoring update", *device.DisplayName)
				return device, nil
			}
			// the clusterName is the same and hostName is not the same, need update
			*device.DisplayName = *existingDevice.DisplayName
			newDevice, err2 := m.updateAndReplace(existingDevice.ID, device)
			if err2 != nil {
				return nil, err2
			}
			log.Infof("Updating existing device (%s)", *newDevice.DisplayName)
			return newDevice, nil
		}
	}
	// duplicate device exists. update displayName and re-add
	renamedDevice, err := m.renameAndAddDevice(device)
	if err != nil {
		log.Errorf("rename device failed: %v", err)
		return nil, fmt.Errorf("rename device failed")
	}
	return renamedDevice, nil
}

// renameAndAddDevice rename display name and then add the device
func (m *Manager) renameAndAddDevice(device *models.Device) (*models.Device, error) {
	resourceName := m.GetPropertyValue(device, constants.K8sResourceNamePropertyKey)
	if resourceName == "" {
		resourceName = *device.DisplayName
	}
	renameResourceName := fmt.Sprintf("%s-%s", resourceName, m.Config().ClusterName)
	existingDevice, err := m.FindByDisplayName(renameResourceName)
	if err != nil {
		log.Warnf("Get device(%s) failed, err: %v", resourceName, err)
	}
	if existingDevice != nil {
		if m.Config().ClusterName == m.GetPropertyValue(existingDevice, constants.K8sClusterNamePropertyKey) {
			device.DisplayName = existingDevice.DisplayName
			return m.updateAndReplace(existingDevice.ID, device)
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
	device := buildDevice(m.Config(), nil, options...)
	log.Debugf("%#v", device)

	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(device)
	restResponse, err := m.LMClient.LM.AddDevice(params)
	if err != nil {
		deviceDefault, ok := err.(*lm.AddDeviceDefault)
		if !ok {
			return nil, err
		}
		// handle the device existing case
		if deviceDefault != nil && deviceDefault.Code() == 409 {
			log.Infof("Check and Update the existing device: %s", *device.DisplayName)
			newDevice, err2 := m.checkAndUpdateExistingDevice(device)
			if err2 != nil {
				return nil, err2
			}
			m.DC.Set(*newDevice.DisplayName)
			return newDevice, nil
		}

		return nil, err
	}
	log.Debugf("%#v", restResponse)
	m.DC.Set(*restResponse.Payload.DisplayName)
	return restResponse.Payload, nil
}

// UpdateAndReplace implements types.DeviceManager.
func (m *Manager) UpdateAndReplace(d *models.Device, options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), d, options...)
	log.Debugf("%#v", device)

	return m.updateAndReplace(d.ID, device)
}

// UpdateAndReplaceByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByDisplayName(name string, filter types.UpdateFilter, options ...types.DeviceOption) (*models.Device, error) {
	if !m.DC.Exists(name) {
		log.Infof("Missing device %v; adding it now", name)
		return m.Add(options...)
	}
	if filter != nil && !filter() {
		log.Debugf("filtered device update %s", name)
		return nil, nil
	}

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
	device, err := m.UpdateAndReplace(d, options...)
	if err != nil {

		return nil, err
	}
	m.DC.Set(*device.DisplayName)
	return device, nil
}

// TODO: this method needs to be removed in DEV-50496

// UpdateAndReplaceField implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceField(d *models.Device, field string, options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), d, options...)
	log.Debugf("%#v", device)

	params := lm.NewPatchDeviceParams()
	params.SetID(d.ID)
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
	device, err := m.UpdateAndReplaceField(d, field, options...)
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
	if err == nil {
		m.DC.Unset(name)
	}
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
