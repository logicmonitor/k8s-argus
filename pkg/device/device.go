package device

import (
	"context"
	"fmt"

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

// checkAndUpdateExistingDevice tries to find and update the devices which needs to be changed
func (m *Manager) checkAndUpdateExistingDevice(device *models.Device) (*models.Device, error) {
	oldDevice, err := m.FindByDisplayName(*device.DisplayName)
	if err != nil {
		return nil, err
	}
	if oldDevice == nil {
		return nil, fmt.Errorf("can not find the device: %s", *device.DisplayName)
	}

	// the device which is not changed will be ignored
	if *device.Name == *oldDevice.Name {
		log.Infof("No changes to device (%s). Ignoring update", *device.DisplayName)
		return device, nil
	}

	// the device of the other cluster will be ignored
	oldClusterName := ""
	if oldDevice.CustomProperties != nil && len(oldDevice.CustomProperties) > 0 {
		for _, cp := range oldDevice.CustomProperties {
			if *cp.Name == constants.K8sClusterNamePropertyKey {
				oldClusterName = *cp.Value
			}
		}
	}
	if oldClusterName != m.Config().ClusterName {
		log.Infof("Device (%s) belongs to a different cluster (%s). Ignoring update", *device.DisplayName, oldClusterName)
		return device, nil
	}

	newDevice, err := m.updateAndReplace(oldDevice.ID, device)
	if err != nil {
		return nil, err
	}
	log.Infof("Finished updating the device: %s", *newDevice.DisplayName)
	return newDevice, nil
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

// Add implements types.DeviceManager.
func (m *Manager) Add(options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), m.ControllerClient, options...)
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
			return newDevice, nil
		}

		return nil, err
	}
	log.Debugf("%#v", restResponse)

	return restResponse.Payload, nil
}

// UpdateAndReplaceByID implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByID(id int32, options ...types.DeviceOption) (*models.Device, error) {
	device := buildDevice(m.Config(), m.ControllerClient, options...)
	log.Debugf("%#v", device)

	return m.updateAndReplace(id, device)
}

// UpdateAndReplaceByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByDisplayName(name string, options ...types.DeviceOption) (*models.Device, error) {
	d, err := m.FindByDisplayName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Warnf("Could not find device %q", name)
		return nil, nil
	}

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
	d, err := m.FindByDisplayName(name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Infof("Could not find device %q", name)
		return nil, nil
	}

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
	d, err := m.FindByDisplayName(name)
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
