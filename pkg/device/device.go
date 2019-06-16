package device

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"

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
	analysisDevices, err := m.FindByDisplayNameContains(*device.DisplayName)
	if err != nil {
		return nil, err
	}
	if len(analysisDevices) == 0 {
		return nil, fmt.Errorf("can not find the devices with string : %s", *device.DisplayName)
	}

	isSame, isUpdate, isAdd, compareDevice := m.analysisDevices(analysisDevices, device)

	// the device which is not changed will be ignored
	if isSame {
		log.Infof("No changes to device (%s). Ignoring update", *device.DisplayName)
		return device, nil
	}

	// the clusterName is the same and hostName is not the same, need update
	if isUpdate {
		newDevice, err := m.updateAndReplace(compareDevice.ID, device)
		if err != nil {
			return nil, err
		}
		log.Infof("Finished updating the device: %s", *newDevice.DisplayName)
		return newDevice, nil
	}

	// the clusterName is not the same, rename displayName and add device
	if isAdd {
		return m.renameAndAddDevice(device)
	}
	return nil, fmt.Errorf("failed to analysis devices: %s", *device.DisplayName)

}

func (m *Manager) analysisDevices(analysisDevices []*models.Device, device *models.Device) (isSame bool, isUpdate bool, isAdd bool, compareDevice *models.Device) {
	if len(analysisDevices) == 0 {
		return
	}
	for _, analysisDevice := range analysisDevices {
		clusterName := m.GetPropertyValue(analysisDevice, constants.K8sClusterNamePropertyKey)
		if clusterName == m.Config().ClusterName {
			if *analysisDevice.Name == *device.Name {
				isSame = true
				compareDevice = analysisDevice
				return isSame, isUpdate, isAdd, compareDevice
			}
			isUpdate = true
			compareDevice = analysisDevice
			return isSame, isUpdate, isAdd, compareDevice
		}
	}
	isAdd = true
	return isSame, isUpdate, isAdd, compareDevice
}

// renameAndAddDevice rename display name and then add the device
func (m *Manager) renameAndAddDevice(device *models.Device) (*models.Device, error) {
	log.Infof("Start rename device(%s)", *device.DisplayName)
	resourceName := m.GetPropertyValue(device, constants.K8sResourceNamePropertyKey)
	if resourceName == "" {
		resourceName = *device.DisplayName
	}
	if resourceName == "" {
		return nil, fmt.Errorf("get device(%s) resource name failed", *device.DisplayName)
	}
	renameResourceName, err := m.getAvailableDisplayName(resourceName)
	if err != nil {
		return nil, err
	}
	if renameResourceName == "" {
		return nil, fmt.Errorf("device(%s) rename failed", *device.DisplayName)
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

// getAvailableDisplayName get the available display name
func (m *Manager) getAvailableDisplayName(resourceName string) (string, error) {
	fields := "displayName"
	queryParams := lm.NewGetDeviceListParams()
	queryParams.Fields = &fields

	startIndex := 0
	renameResourceName := ""
	for startIndex < 10 {
		displayNames := utilities.GetBatchDisplayNames(resourceName, m.Config().ClusterName, 5, &startIndex)
		log.Debugf("Get batch display names: %+v", displayNames)
		if len(displayNames) == 0 {
			return renameResourceName, errors.New("can't get batch display names")
		}
		filter := fmt.Sprintf("displayName:\"%s\"", strings.Join(displayNames, "\"|\""))
		queryParams.SetFilter(&filter)
		deviceList, err := m.LMClient.LM.GetDeviceList(queryParams)
		if err != nil {
			return renameResourceName, err
		}
		log.Debugf("Get device list by displayNames: %+v", deviceList.Payload.Items)
		for _, dn := range displayNames {
			flag := true
			for _, item := range deviceList.Payload.Items {
				if dn == *item.DisplayName {
					flag = false
					break
				}
			}
			if flag {
				renameResourceName = dn
				break
			}
		}
		if renameResourceName != "" {
			break
		}
	}
	return renameResourceName, nil
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

// FindByDisplayNameContains implements types.DeviceManager.
func (m *Manager) FindByDisplayNameContains(name string) ([]*models.Device, error) {
	filter := fmt.Sprintf("displayName~\"%s\"", name)
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, err
	}
	log.Debugf("%#v", restResponse)
	return restResponse.Payload.Items, nil
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
