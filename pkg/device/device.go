package device

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device/builder"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"

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
	displayNameWithClusterName := fmt.Sprintf("%s-%s", *device.DisplayName, m.Config().ClusterName)
	existingDevices, err := m.FindByDisplayNames(lctx, resource, *device.DisplayName, displayNameWithClusterName)
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
			newDevice, err2 := m.updateAndReplace(lctx, resource, existingDevice.ID, device)
			if err2 != nil {
				return nil, err2
			}
			log.Infof("Updating existing device (%s)", *newDevice.DisplayName)
			return newDevice, nil
		}
	}
	// duplicate device exists. update displayName and re-add
	renamedDevice, err := m.renameAndAddDevice(lctx, resource, device)
	if err != nil {
		log.Errorf("rename device failed: %v", err)
		return nil, fmt.Errorf("rename device failed")
	}
	return renamedDevice, nil
}

// renameAndAddDevice rename display name and then add the device
func (m *Manager) renameAndAddDevice(lctx *lmctx.LMContext, resource string, device *models.Device) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	resourceName := m.GetPropertyValue(device, constants.K8sResourceNamePropertyKey)
	if resourceName == "" {
		resourceName = *device.DisplayName
	}
	renameResourceName := fmt.Sprintf("%s-%s", resourceName, m.Config().ClusterName)
	existingDevice, err := m.FindByDisplayName(lctx, resource, renameResourceName)
	if err != nil {
		log.Warnf("Get device(%s) failed, err: %v", resourceName, err)
	}
	if existingDevice != nil {
		if m.Config().ClusterName == m.GetPropertyValue(existingDevice, constants.K8sClusterNamePropertyKey) {
			device.DisplayName = existingDevice.DisplayName
			return m.updateAndReplace(lctx, resource, existingDevice.ID, device)
		}
		return nil, fmt.Errorf("exist displayName: %s", renameResourceName)
	}
	log.Infof("Rename device: %s -> %s", *device.DisplayName, renameResourceName)
	device.DisplayName = &renameResourceName
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
	//restResponse, err := m.LMClient.LM.AddDevice(params)
	if err != nil {
		return nil, err
	}
	return restResponse.(*lm.AddDeviceOK).Payload, nil
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

// FindByDisplayNameAndClusterName implements types.DeviceManager.
func (m *Manager) FindByDisplayNameAndClusterName(lctx *lmctx.LMContext, resource string, displayName string) (*models.Device, error) {
	displayNameWithClusterName := fmt.Sprintf("%s-%s", displayName, m.Config().ClusterName)
	devices, err := m.FindByDisplayNames(lctx, resource, displayName, displayNameWithClusterName)
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

// getEvaluationParamsForResource generates evaluation parameters based on labels and specified resource
func getEvaluationParamsForResource(device *models.Device, labels map[string]string) (map[string]interface{}, error) {
	evaluationParams := make(map[string]interface{})

	for key, value := range labels {
		key = filters.CheckAndReplaceInvalidChars(key)
		value = filters.CheckAndReplaceInvalidChars(value)
		evaluationParams[key] = value
	}

	evaluationParams["name"] = *device.DisplayName
	return evaluationParams, nil
}

// Add implements types.DeviceManager.
func (m *Manager) Add(lctx *lmctx.LMContext, resource string, labels map[string]string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), nil, options...)
	log.Debugf("%#v", device)

	evaluationParams, err := getEvaluationParamsForResource(device, labels)
	if err != nil {
		return nil, err
	}
	log.Debugf("Evaluation params for resource %s %+v:", resource, evaluationParams)

	if filters.Eval(resource, evaluationParams) {
		log.Infof("Filtering out %s %s.", resource, *device.DisplayName)
		// delete existing resource which is mentioned for filtering.
		err := m.DeleteByDisplayName(lctx, resource, *device.DisplayName)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

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
	if err != nil {
		deviceDefault, ok := err.(*lm.AddDeviceDefault)
		if !ok {
			return nil, err
		}
		// handle the device existing case
		if deviceDefault != nil && deviceDefault.Code() == 409 {
			log.Infof("Check and Update the existing device: %s", *device.DisplayName)
			newDevice, err2 := m.checkAndUpdateExistingDevice(lctx, resource, device)
			if err2 != nil {
				return nil, err2
			}
			m.DC.Set(*newDevice.DisplayName)
			return newDevice, nil
		}

		return nil, err
	}
	resp := restResponse.(*lm.AddDeviceOK)
	m.DC.Set(*resp.Payload.DisplayName)
	log.Debugf("%#v", resp)
	return resp.Payload, nil
}

// UpdateAndReplace implements types.DeviceManager.
func (m *Manager) UpdateAndReplace(lctx *lmctx.LMContext, resource string, d *models.Device, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), d, options...)
	log.Debugf("%#v", device)

	return m.updateAndReplace(lctx, resource, d.ID, device)
}

// UpdateAndReplaceByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceByDisplayName(lctx *lmctx.LMContext, resource string, name string, filter types.UpdateFilter, labels map[string]string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	if !m.DC.Exists(name) {
		log.Infof("Missing device %v; adding it now", name)
		return m.Add(lctx, resource, labels, options...)
	}
	if filter != nil && !filter() {
		log.Debugf("filtered device update %s", name)
		return nil, nil
	}

	d, err := m.FindByDisplayNameAndClusterName(lctx, resource, name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Warnf("Could not find device %q", name)
		return nil, nil
	}

	options = append(options, m.DisplayName(*d.DisplayName))
	// Update the device.
	device, err := m.UpdateAndReplace(lctx, resource, d, options...)
	if err != nil {

		return nil, err
	}
	m.DC.Set(*device.DisplayName)
	return device, nil
}

// TODO: this method needs to be removed in DEV-50496

// UpdateAndReplaceField implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceField(lctx *lmctx.LMContext, resource string, d *models.Device, field string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	device := buildDevice(lctx, m.Config(), d, options...)
	log.Debugf("%#v", device)

	params := lm.NewPatchDeviceParams()
	params.SetID(d.ID)
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
	//restResponse, err := m.LMClient.LM.PatchDevice(params)
	if err != nil {
		return nil, err
	}
	resp := restResponse.(*lm.PatchDeviceOK)
	log.Debugf("%#v", resp)

	return resp.Payload, nil
}

// TODO: this method needs to be removed in DEV-50496

// UpdateAndReplaceFieldByDisplayName implements types.DeviceManager.
func (m *Manager) UpdateAndReplaceFieldByDisplayName(lctx *lmctx.LMContext, resource string, name string, field string, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	d, err := m.FindByDisplayNameAndClusterName(lctx, resource, name)
	if err != nil {
		return nil, err
	}

	if d == nil {
		log.Infof("Could not find device %q", name)
		return nil, nil
	}
	options = append(options, m.DisplayName(*d.DisplayName))
	// Update the device.
	device, err := m.UpdateAndReplaceField(lctx, resource, d, field, options...)
	if err != nil {
		return nil, err
	}

	return device, nil
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
func (m *Manager) DeleteByDisplayName(lctx *lmctx.LMContext, resource string, name string) error {
	log := lmlog.Logger(lctx)
	d, err := m.FindByDisplayNameAndClusterName(lctx, resource, name)
	if err != nil {
		return err
	}

	// TODO: Should this return an error?
	if d == nil {
		log.Infof("Could not find device %q", name)
		return nil
	}
	err2 := m.DeleteByID(lctx, resource, d.ID)
	if err2 == nil {
		m.DC.Unset(name)
		log.Infof("deleted device %q", name)
	}
	return err2
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
	fields := "id,name,displayName,customProperties"
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
