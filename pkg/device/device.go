package device

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"github.com/logicmonitor/k8s-argus/pkg/lmexec"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	cscutils "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// Manager implements types.DeviceManager
type Manager struct {
	ResourceType string
	*types.Base
	*builder.Builder
	types.LMExecutor
	types.LMFacade
	DC *devicecache.DeviceCache
}

func buildDevice(lctx *lmctx.LMContext, c *config.Config, d *models.Device, options ...types.DeviceOption) *models.Device {
	log := lmlog.Logger(lctx)
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
func (m *Manager) checkAndUpdateExistingDevice(lctx *lmctx.LMContext, resource string, device *models.Device) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	currentCluster := m.Config().ClusterName
	existingDevices, err := m.FindByDisplayNames(lctx, resource, *device.DisplayName, util.GetFullDisplayName(device, resource, currentCluster))

	if err != nil {
		return nil, err
	}
	if len(existingDevices) == 0 {
		return nil, fmt.Errorf("could not find devices with names: %s , %s", *device.DisplayName, util.GetFullDisplayName(device, resource, currentCluster))
	}
	for _, existingDevice := range existingDevices {
		clusterName := util.GetPropertyValue(existingDevice, constants.K8sClusterNamePropertyKey)
		if clusterName == currentCluster {
			// the device which is not changed will be ignored
			if util.GetDisplayNameWithNamespace(existingDevice, resource) == util.GetDisplayNameWithNamespace(device, resource) {
				log.Infof("No changes to device (%s). Ignoring update", *device.DisplayName)
				m.DC.Set(util.GetFullDisplayName(existingDevice, resource, currentCluster))
				return nil, nil
			}

			log.Infof("Updating and moving to conflicts group - existing device (%s)", *existingDevice.DisplayName)
			err2 := m.moveDeviceToConflictGroup(lctx, existingDevice, resource)
			if err2 != nil {
				log.Errorf("%v", err2)
				return nil, err2
			}

			m.DC.Set(util.GetFullDisplayName(existingDevice, resource, currentCluster))
			return existingDevice, nil
		}
	}
	return device, nil
}

// renameAndAddDevice rename display name and then add the device
func (m *Manager) renameAndAddDevice(lctx *lmctx.LMContext, resource string, device *models.Device) (*models.Device, error) {
	restResponse, err := m.addDevice(lctx, resource, device)

	if err != nil {
		return nil, err
	}
	return restResponse.(*lm.AddDeviceOK).Payload, nil
}

// UpdateDeviceName updates the device hostname
func (m *Manager) UpdateDeviceName(lctx *lmctx.LMContext, resource string, device *models.Device, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	newDevice := buildDevice(lctx, m.Config(), device, options...)

	fields := constants.NameFieldName
	updatedDevice, err := m.UpdateAndReplaceField(lctx, resource, newDevice, fields)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	m.DC.Set(util.GetFullDisplayName(updatedDevice, resource, m.Config().ClusterName))
	return updatedDevice, nil
}

// RenameAndUpdateDevice renames the device display as per desiredDisplayName and moves the conflicting devices to conflicts dynamic group.
func (m *Manager) RenameAndUpdateDevice(lctx *lmctx.LMContext, resource string, device *models.Device, desiredDisplayName string) error {
	log := lmlog.Logger(lctx)
	collectorID := cscutils.GetCollectorID()
	device.PreferredCollectorID = &collectorID

	*device.DisplayName = desiredDisplayName
	updatedDevice, err := m.updateAndReplace(lctx, resource, device.ID, device)

	// remove conflict category from sys-category of device.
	if util.IsConflictingDevice(updatedDevice, resource) {
		currentCategories := getCurrentSystemCategoriesForDevice(updatedDevice)
		updatedCategories := strings.ReplaceAll(currentCategories, util.GetConflictCategoryByResourceType(resource), "")
		entityProperty := models.EntityProperty{Name: constants.K8sSystemCategoriesPropertyKey, Value: updatedCategories, Type: "system"}
		err1 := m.updateDevicePropertyByName(lctx, updatedDevice.ID, &entityProperty, resource)
		if err1 != nil {
			log.Errorf("failed to remove device %s from conflicts group. %v", *updatedDevice.DisplayName, err1)
			return err1
		}
	}

	if err != nil {
		deviceDefault, _ := err.(*lm.UpdateDeviceDefault)
		// handle the device existing case
		if deviceDefault != nil && deviceDefault.Code() == 409 {
			log.Infof("Device with displayName %s already exists, rename and moving it to conflicts group.", *device.DisplayName)
			options := []types.DeviceOption{
				m.DisplayName(util.GetFullDisplayName(device, resource, m.Config().ClusterName)),
			}
			newDevice, err := m.UpdateAndReplace(lctx, resource, device, options...)
			if err != nil {
				log.Errorf("failed to update the device %s : %v", *device.DisplayName, err)
				return err
			}
			err2 := m.moveDeviceToConflictGroup(lctx, newDevice, resource)
			if err2 != nil {
				log.Errorf("failed to move device %s to conflicts group: %v", *newDevice.DisplayName, err2)
				return err2
			}

			m.DC.Set(util.GetFullDisplayName(newDevice, resource, m.Config().ClusterName))
			return nil
		}
		log.Errorf("%v", err)
		return err
	}

	m.DC.Set(util.GetFullDisplayName(updatedDevice, resource, m.Config().ClusterName))
	return nil
}

// GetDesiredDisplayName returns desired display name based on FullDisplayNameIncludeClusterName and FullDisplayNameIncludeNamespace properties.
func (m *Manager) GetDesiredDisplayName(name, namespace, resource string) string {
	return util.GetDesiredDisplayNameByResourceAndConfig(name, namespace, m.Config().ClusterName, resource, m.Config().FullDisplayNameIncludeNamespace, m.Config().FullDisplayNameIncludeClusterName)
}

func (m *Manager) updateAndReplace(lctx *lmctx.LMContext, resource string, id int32, device *models.Device) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	opType := "replace"
	params := lm.NewUpdateDeviceParams()
	params.SetID(id)
	params.SetBody(device)
	params.SetOpType(&opType)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.UpdateDevice(params),
			LMCtx:   lctx,
		},
		Method:   "PUT",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.UpdateDeviceErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)

	//restResponse, err := m.LMClient.LM.UpdateDevice(params)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.UpdateDeviceOK)
	log.Debugf("%#v", resp)

	return resp.Payload, nil
}

// FindByDisplayName implements types.DeviceManager.
func (m *Manager) FindByDisplayName(lctx *lmctx.LMContext, resource string, name string) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	filter := fmt.Sprintf("displayName:\"%s\"", name)
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.GetDeviceList(params),
			LMCtx:   lctx,
		},
		Method:   "GET",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.GetDeviceListErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	//restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.GetDeviceListOK)
	log.Debugf("%#v", resp)
	if resp.Payload.Total == 1 {
		return resp.Payload.Items[0], nil
	}

	return nil, nil
}

// FindByDisplayNames implements types.DeviceManager.
func (m *Manager) FindByDisplayNames(lctx *lmctx.LMContext, resource string, displayNames ...string) ([]*models.Device, error) {
	log := lmlog.Logger(lctx)
	if len(displayNames) == 0 {
		return []*models.Device{}, nil
	}
	filter := fmt.Sprintf("displayName:\"%s\"", strings.Join(displayNames, "\"|\""))
	params := lm.NewGetDeviceListParams()
	params.SetFilter(&filter)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.GetDeviceList(params),
			LMCtx:   lctx,
		},
		Method:   "GET",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.GetDeviceListErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	//restResponse, err := m.LMClient.LM.GetDeviceList(params)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.GetDeviceListOK)
	log.Debugf("%#v", resp)
	return resp.Payload.Items, nil
}

// getEvaluationParamsForResource generates evaluation parameters based on labels and specified resource
func getEvaluationParamsForResource(device *models.Device, labels map[string]string) map[string]interface{} {
	evaluationParams := make(map[string]interface{})

	for key, value := range labels {
		key = filters.CheckAndReplaceInvalidChars(key)
		value = filters.CheckAndReplaceInvalidChars(value)
		evaluationParams[key] = value
	}

	evaluationParams["name"] = filters.CheckAndReplaceInvalidChars(*device.DisplayName)
	return evaluationParams
}

func (m *Manager) evalExclusion(lctx *lmctx.LMContext, device *models.Device, resource string, labels map[string]string) bool {
	log := lmlog.Logger(lctx)
	evaluationParams := getEvaluationParamsForResource(device, labels)
	log.Debugf("Evaluation params for resource %s %+v:", resource, evaluationParams)
	return filters.Eval(resource, evaluationParams)
}

// Add implements types.DeviceManager.
func (m *Manager) Add(lctx *lmctx.LMContext, resource string, labels map[string]string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), nil, options...)

	if !m.checkPingDeviceAndSystemIPs(lctx, device) {
		log.Warnf("Property '%s' is empty for device '%s', skipping", constants.K8sSystemIPsPropertyKey, *device.DisplayName)
		return nil, nil
	}

	isExclude := m.evalExclusion(lctx, device, resource, labels)

	if isExclude {
		// delete existing resource which is mentioned for filtering.
		log.Infof("Device %s is being excluded, sending delete if device present", *device.DisplayName)
		err := m.DeleteByDisplayName(lctx, resource, *device.DisplayName, util.GetFullDisplayName(device, resource, m.Config().ClusterName))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	restResponse, err := m.addDevice(lctx, resource, device)
	if err != nil {
		deviceDefault, ok := err.(*lm.AddDeviceDefault)
		if !ok {
			return nil, err
		}
		// handle the device existing case
		if deviceDefault != nil && deviceDefault.Code() == 409 {
			newdevice, err := m.addConflictingDevice(lctx, device, resource, options...)
			if err != nil {
				return nil, err
			}
			return newdevice, nil
		}
		return nil, err
	}
	resp := restResponse.(*lm.AddDeviceOK)
	m.DC.Set(util.GetFullDisplayName(resp.Payload, resource, m.Config().ClusterName))
	log.Debugf("%#v", resp)
	return resp.Payload, nil
}

func (m *Manager) addConflictingDevice(lctx *lmctx.LMContext, device *models.Device, resource string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	log.Infof("Check and Update the existing device: %s", *device.DisplayName)
	updatedevice, err := m.checkAndUpdateExistingDevice(lctx, resource, device)
	if err != nil {
		err2 := fmt.Errorf("failed to update existing device: %w", err)
		log.Error(err2.Error())
		return nil, err2
	}

	if updatedevice == nil {
		return device, nil
	}

	currentCluster := m.Config().ClusterName
	log.Infof("Adding new device %s and moving to conflicts group.", *device.DisplayName)
	options = append(options, m.Custom(constants.K8sSystemCategoriesPropertyKey, util.GetConflictCategoryByResourceType(resource)))
	*device.DisplayName = util.GetFullDisplayName(device, resource, currentCluster)
	newDevice := buildDevice(lctx, m.Config(), device, options...)
	renamedDevice, err := m.renameAndAddDevice(lctx, resource, newDevice)

	if err != nil {
		log.Errorf("add new device failed: %v", err)
		return nil, err
	}

	m.DC.Set(util.GetFullDisplayName(renamedDevice, resource, currentCluster))
	return renamedDevice, nil
}

func (m *Manager) moveDeviceToConflictGroup(lctx *lmctx.LMContext, device *models.Device, resource string) error {
	updatedCategories := fmt.Sprintf(getCurrentSystemCategoriesForDevice(device) + "," + util.GetConflictCategoryByResourceType(resource))
	entityProperty := models.EntityProperty{Name: constants.K8sSystemCategoriesPropertyKey, Value: updatedCategories, Type: "system"}
	err := m.updateDevicePropertyByName(lctx, device.ID, &entityProperty, resource)
	if err != nil {
		return err
	}
	return nil
}

func getCurrentSystemCategoriesForDevice(device *models.Device) string {
	return util.GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)
}

// checkPingDeviceAndSystemIPs verifies that 'system.ips' property is present if device ping feature is enabled.
// If hostNetwork is enabled then device hostname is set as resource name instead of IP Address.
// In this case collector uses 'system.ips' to communicate with the resource.
func (m *Manager) checkPingDeviceAndSystemIPs(lctx *lmctx.LMContext, device *models.Device) bool {
	isPingDevice := lctx.Extract(constants.IsPingDevice)
	if isPingDevice != nil && isPingDevice.(bool) && util.GetPropertyValue(device, constants.K8sSystemIPsPropertyKey) == "" {
		return false
	}
	return true
}

func (m *Manager) addDevice(lctx *lmctx.LMContext, resource string, device *models.Device) (interface{}, error) {
	params := lm.NewAddDeviceParams()
	addFromWizard := false
	params.SetAddFromWizard(&addFromWizard)
	params.SetBody(device)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.AddDevice(params),
			LMCtx:   lctx,
		},
		Method:   "POST",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.AddDeviceErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	return restResponse, err
}

// UpdateAndReplace implements types.DeviceManager.
func (m *Manager) UpdateAndReplace(lctx *lmctx.LMContext, resource string, d *models.Device, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), d, options...)
	log.Debugf("%#v", device)

	return m.updateAndReplace(lctx, resource, d.ID, device)
}

// UpdateAndReplaceByDisplayName implements types.DeviceManager.
// nolint
func (m *Manager) UpdateAndReplaceByDisplayName(lctx *lmctx.LMContext, resource, name, fullName string, filter types.UpdateFilter, labels map[string]string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), nil, options...)

	if *device.Name == "" {
		err := fmt.Errorf("empty name is invalid for resource %v", resource)
		log.Warn(err.Error())
		return nil, err
	}

	isExclude := m.evalExclusion(lctx, device, resource, labels)
	if isExclude {
		// delete existing resource which is mentioned for discovery filtering.
		log.Infof("Device %s is excluded, sending delete if device present", *device.DisplayName)
		if m.DC.Exists(fullName) {
			err := m.DeleteByDisplayName(lctx, resource, *device.DisplayName, util.GetFullDisplayName(device, resource, m.Config().ClusterName))
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	if !m.DC.Exists(fullName) {
		log.Infof("Missing device %v; (full name = %v) adding it now", name, fullName)
		return m.Add(lctx, resource, labels, options...)
	}

	if filter != nil && !filter() {
		log.Debugf("Ignoring updates for device %s, %s", name, resource)
		return nil, nil
	}

	existingDevice, err := m.getExisitingDeviceByGivenProperties(lctx, name, fullName, resource)

	if err != nil {
		return nil, err
	}

	if existingDevice == nil {
		log.Warnf("Could not find device %q", name)
		return nil, nil
	}

	if resource == constants.Pods && filter != nil && filter() {
		log.Warnf("IP address mismatch for device %s", *device.DisplayName)
		existingDevice, err = m.UpdateDeviceName(lctx, resource, existingDevice, append(options, m.Name(util.GetPropertyValue(device, constants.K8sSystemIPsPropertyKey)))...)
		if util.GetPropertyValue(device, constants.K8sSystemIPsPropertyKey) != *device.Name {
			err := m.WaitToUpdateSysIps(lctx, existingDevice, resource, util.GetPropertyValue(device, constants.K8sSystemIPsPropertyKey), 1*time.Minute)
			if err != nil {
				return nil, err
			}
		}
	}

	options = append(options, m.DisplayName(*existingDevice.DisplayName))

	// Update the device.
	updatedDevice, err := m.UpdateAndReplace(lctx, resource, existingDevice, options...)
	if err != nil {

		return nil, err
	}
	m.DC.Set(util.GetFullDisplayName(updatedDevice, resource, m.Config().ClusterName))
	return device, nil
}

// WaitToUpdateSysIps waits until system.ips property is updated.
func (m *Manager) WaitToUpdateSysIps(lctx *lmctx.LMContext, device *models.Device, resourceType string, address string, timeout time.Duration) error {
	backOff := 5 * time.Second
	c, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-c.Done():
			return fmt.Errorf("timed out (5 minutes) waiting to change system.ips property on device: %s", *device.DisplayName)

		default:
			isIPUpdated, err := m.VerifyIfPodIPUpdated(lctx, device, resourceType, address)
			if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
				return fmt.Errorf("failed to update IP address for %s from k8s, err: %w", *device.DisplayName, err)
			}
			if isIPUpdated {
				return nil
			}
			time.Sleep(backOff)
			backOff *= 2
		}
	}
}

// VerifyIfPodIPUpdated checks if system.ips property is updated
func (m *Manager) VerifyIfPodIPUpdated(lctx *lmctx.LMContext, device *models.Device, resourceType, expectedPodIP string) (bool, error) {
	log := lmlog.Logger(lctx)
	existingDevice, err := m.FindByDisplayName(lctx, strings.ToLower(resourceType), *device.DisplayName)

	if err != nil {
		return false, err
	}

	podIP := util.GetPropertyValue(existingDevice, constants.K8sSystemIPsPropertyKey)
	log.Debugf("current ip = %s and expected IP = %s", podIP, expectedPodIP)
	return podIP == expectedPodIP, err
}

// UpdateAndReplaceField implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceField(lctx *lmctx.LMContext, resource string, device *models.Device, fields string) (*models.Device, error) {
	params := lm.NewPatchDeviceParams()
	params.SetID(device.ID)
	params.SetPatchFields(&fields)
	params.SetBody(device)
	opType := "replace"
	params.SetOpType(&opType)

	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.PatchDevice(params),
			LMCtx:   lctx,
		},
		Method:   "PATCH",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.PatchDeviceErrResp,
		},
	}

	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.PatchDeviceOK)
	return resp.Payload, nil
}

// UpdateDevicePropertyByName updates the specified property value for a device.
func (m *Manager) updateDevicePropertyByName(lctx *lmctx.LMContext, deviceID int32, entityProperty *models.EntityProperty, resource string) error {
	log := lmlog.Logger(lctx)
	params := lm.NewUpdateDevicePropertyByNameParams()
	params.SetBody(entityProperty)
	params.SetDeviceID(deviceID)
	params.SetName(entityProperty.Name)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.UpdateDevicePropertyByName(params),
			LMCtx:   lctx,
		},
		Method:   "PUT",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.UpdateDevicePropertyErrResp,
		},
	}
	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	//restResponse, err := client.LM.UpdateDevicePropertyByName(params)

	if err != nil {
		return fmt.Errorf("failed to update device property '%v'. Error: %v", entityProperty.Name, err)
	}
	resp := restResponse.(*lm.UpdateDevicePropertyByNameOK)
	log.Debugf("Update property response payload : %#v", resp.Payload)

	return nil
}

// MoveToDeletedGroup implements types.DeviceManager.
func (m *Manager) MoveToDeletedGroup(lctx *lmctx.LMContext, resource, name, fullName string, deletionTimestamp *v1.Time, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	existingDevice, err := m.getExisitingDeviceByGivenProperties(lctx, name, fullName, resource)
	if err != nil {
		return nil, err
	}
	if existingDevice == nil {
		log.Warnf("Could not find device %q", name)
		return nil, nil
	}

	device := m.buildDeviceBeforeDeletion(deletionTimestamp, existingDevice, options...)
	fields := constants.CustomPropertiesFieldName + "," + constants.NameFieldName + "," + constants.DisplayNameFieldName

	updatedDevice, err := m.UpdateAndReplaceField(lctx, resource, device, fields)
	if err != nil {
		return nil, err
	}
	return updatedDevice, nil
}

func (m *Manager) buildDeviceBeforeDeletion(deletionTimestamp *v1.Time, existingDevice *models.Device, options ...types.DeviceOption) *models.Device {
	// add resource deletion timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if deletionTimestamp != nil {
		timestamp = strconv.FormatInt(deletionTimestamp.Unix(), 10)
	}
	options = append(options, m.DeletedOn(constants.K8sResourceDeletedOnPropertyKey, timestamp))

	// rename device to resolve conflicts for new devices
	shortUUID := strconv.FormatUint(uint64(util.GetShortUUID()), 10)
	deviceName := util.TrimName(*existingDevice.Name)
	options = append(options, m.Name(deviceName+"-"+shortUUID))
	deviceDisplayName := util.TrimName(util.GetPropertyValue(existingDevice, constants.K8sDeviceNamePropertyKey))
	options = append(options, m.DisplayName(deviceDisplayName+"-"+shortUUID))

	// build device with specific fields that needs to be updated
	// ID & PreferredCollectorID are required, if not passed then considered as 0 & API throws an error
	device := &models.Device{
		ID:                   existingDevice.ID,
		Name:                 existingDevice.Name,
		DisplayName:          existingDevice.DisplayName,
		CustomProperties:     existingDevice.CustomProperties,
		PreferredCollectorID: existingDevice.PreferredCollectorID,
	}

	for _, option := range options {
		option(device)
	}

	return device
}

// DeleteByID implements types.DeviceManager.
func (m *Manager) DeleteByID(lctx *lmctx.LMContext, resource string, id int32) error {
	params := lm.NewDeleteDeviceByIDParams()
	params.SetID(id)
	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.DeleteDeviceByID(params),
			LMCtx:   lctx,
		},
		Method:   "DELETE",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.DeleteDeviceByIDErrResp,
		},
	}
	_, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	//_, err := m.LMClient.LM.DeleteDeviceByID(params)
	return err
}

// DeleteByDisplayName implements types.DeviceManager.
func (m *Manager) DeleteByDisplayName(lctx *lmctx.LMContext, resource, name, fullName string) error {
	log := lmlog.Logger(lctx)
	existingDevice, err := m.getExisitingDeviceByGivenProperties(lctx, name, fullName, resource)

	if err != nil {
		return err
	}

	if existingDevice == nil {
		log.Infof("Could not find device %q", name)
		m.DC.Unset(fullName)
		return nil
	}

	err2 := m.DeleteByID(lctx, resource, existingDevice.ID)
	if err2 != nil {
		return err2
	}
	m.DC.Unset(fullName)
	log.Infof("Deleted device %q", name)

	return nil
}

func (m *Manager) getExisitingDeviceByGivenProperties(lctx *lmctx.LMContext, name, fullName, resource string) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	existingDevices, err := m.FindByDisplayNames(lctx, resource, name, fullName)

	if err != nil {
		log.Errorf("error finding devices with names: %s, %s - %v", name, fullName, err)
		return nil, err
	}

	if len(existingDevices) == 0 {
		log.Debugf("Could not find device %q", name)
		return nil, nil
	}

	for _, existingDevice := range existingDevices {
		clusterName := util.GetPropertyValue(existingDevice, constants.K8sClusterNamePropertyKey)
		if util.GetFullDisplayName(existingDevice, resource, clusterName) == fullName {
			return existingDevice, nil
		}
	}
	return nil, nil
}

// Config implements types.DeviceManager.
func (m *Manager) Config() *config.Config {
	return m.Base.Config
}

// GetListByGroupID implements getting all the devices belongs to the group directly
func (m *Manager) GetListByGroupID(lctx *lmctx.LMContext, resource string, groupID int32) ([]*models.Device, error) {
	log := lmlog.Logger(lctx)
	params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
	params.SetID(groupID)
	fields := "id,name,displayName,customProperties,systemProperties,inheritedProperties,preferredCollectorId"
	params.SetFields(&fields)
	size := int32(-1)
	params.SetSize(&size)

	cmd := &types.HTTPCommand{
		Command: &types.Command{
			ExecFun: m.GetImmediateDeviceListByDeviceGroupID(params),
			LMCtx:   lctx,
		},
		Method:   "GET",
		Category: "device",
		LMHCErrParse: &types.LMHCErrParse{
			ParseErrResp: m.GetImmediateDeviceListByDeviceGroupIDErrResp,
		},
	}

	restResponse, err := m.LMFacade.SendReceive(lctx, resource, cmd)
	//restResponse, err := m.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.GetImmediateDeviceListByDeviceGroupIDOK)
	log.Debugf("%#v", resp)
	return resp.Payload.Items, nil
}
