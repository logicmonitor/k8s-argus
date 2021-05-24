package device

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// Add implements types.DeviceManager.
func (m *Manager) Add(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return nil, err
	}
	device, err := util.BuildDevice(lctx, conf, nil, options...)
	if err != nil {
		return nil, fmt.Errorf("could not build device object: %w", err)
	}

	objectMeta := *rt.ObjectMeta(obj)
	warn, valid, err2 := validateNewDevice(lctx, rt, device, objectMeta)
	if !valid {
		return nil, err2
	}
	// if warning, then just logging and proceeding further
	if warn {
		log.Warnf("device validation error %s", err2)
	}

	log.Infof("Does device exists in cache")

	_, ok := m.DoesDeviceExistInCache(lctx, rt, device)
	if ok {
		return nil, &types.DeviceExists{}
	}

	log.Infof("Does device conflicts in cache")

	// All conflicts within cluster will be handled here only
	// If conflicts could not be found due to cache miss, then it should be avoided
	_, conflicts := m.DoesDeviceConflictInCluster(lctx, rt, device)
	if conflicts {
		log.Infof("Conflicts with in-cluster resource")
		options = append(options, m.ModifyToUnique(lctx, rt, device, obj)...)
	}

	device, err = util.BuildDevice(lctx, conf, nil, options...)
	if err != nil {
		return nil, fmt.Errorf("could not build device object: %w", err)
	}
	resultDevice, err3 := m.addDevice(lctx, rt, device)
	m2, done, err4 := m.handleConflict(lctx, rt, obj, err3, device)
	if done {
		return m2, err4
	}

	return resultDevice, nil
}

func (m *Manager) handleConflict(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, err3 error, device *models.Device) (*models.Device, bool, error) {
	log := lmlog.Logger(lctx)
	if err3 != nil {
		deviceDefault := err3.(*lm.AddDeviceDefault) // nolint: errorlint
		// Conflicts due to cache miss or outside cluster same device
		// 99.99% call should not reach here: more the control here, more the portal requests
		if deviceDefault != nil && deviceDefault.Code() == http.StatusConflict {
			// nolint: godox
			// TODO: PROM_METRIC cache miss metrics
			log.Infof("Outside cluster conflict")
			options := m.ModifyToUnique(lctx, rt, device, obj)
			conf, err5 := config.GetConfig()
			if err5 != nil {
				return nil, false, err5
			}
			device, err6 := util.BuildDevice(lctx, conf, device, options...)
			if err6 != nil {
				log.Errorf("Failed to buid resource object")
				return nil, false, err6
			}
			resultDevice2, err4 := m.addDevice(lctx, rt, device)
			if err4 != nil {
				log.Errorf("failed to add out-cluster resource: %s", err4)
				// nolint: godox
				// TODO: PROM_METRIC Add device failed metric
				return nil, true, err4
			}

			return resultDevice2, true, nil
		}
		// nolint: godox
		// TODO: PROM_METRIC Add device failed metric

		return nil, true, err3
	}

	return nil, false, nil
}

// Update update
func (m *Manager) Update(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj interface{}, newObj interface{}, options ...types.DeviceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return nil, err
	}
	device, err := util.BuildDevice(lctx, conf, nil, options...)
	if err != nil {
		return nil, err
	}
	log.Tracef("resource built from options: %s", spew.Sdump(device))

	ce, ok := m.DoesDeviceExistInCache(lctx, rt, device)
	if !ok {
		log.Debugf("Device does not exist in cache")

		return nil, fmt.Errorf("device does not exist in cache")
	}
	// If the LM resource is modified outside of argus:
	// For ex:
	// 1. User manually added device property on LM resource 2. PropertySources added property on LM resource
	// Partial patch API is not allowed - Either it will replace all set of customProperties or nothing
	// Similarly, system.categories is comma separated string, so to add or modify, we need to fetch current string from LM resource

	device, err = m.FetchDevice(lctx, rt, ce.LMID)
	if err != nil {
		// This indicates that the lm resource id present in cache is wrong, consider as cache miss
		// nolint: godox
		// TODO: PROM_METRIC cache miss gauge
		log.Errorf("Failed to fetch existing resource: %s", err)

		return nil, err
	}

	_, conflicts := m.DoesDeviceConflictInCluster(lctx, rt, device)
	if conflicts {
		log.Info("Conflicts with in-cluster resource")
		options = append(options, m.ModifyToUnique(lctx, rt, device, newObj)...)
	}
	// Put ID for update request
	device.ID = ce.LMID
	// retain displayName
	options = append(options, m.DisplayName(*device.DisplayName))
	device, err = util.BuildDevice(lctx, conf, device, options...)
	if err != nil {
		return nil, err
	}

	log.Tracef("Updating resource with: %s", spew.Sdump(device))
	// Update the device
	return m.UpdateAndReplaceResource(lctx, rt, device.ID, device)
}

// Delete delete
func (m *Manager) Delete(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.DeviceOption) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return err
	}
	device, err := util.BuildDevice(lctx, conf, nil, options...)
	if err != nil {
		return err
	}
	ce, ok := m.DoesDeviceExistInCache(lctx, rt, device)

	if !ok {
		return nil
	}
	if ce.LMID == 0 {
		m.UnsetDeviceInCache(lctx, rt, device)

		return nil
	}
	device.ID = ce.LMID
	err = m.deleteDevice(lctx, rt, device)
	if err != nil {
		return err
	}

	return nil
}

// MarkDeleted mark
func (m *Manager) MarkDeleted(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.DeviceOption) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return err
	}
	device, err := util.BuildDevice(lctx, conf, nil, options...)
	if err != nil {
		return err
	}
	ce, ok := m.DoesDeviceExistInCache(lctx, rt, device)
	if !ok {
		return nil
	}
	if ce.LMID == 0 {
		m.UnsetDeviceInCache(lctx, rt, device)

		return nil
	}
	device.ID = ce.LMID
	conf, err2 := config.GetConfig()
	if err2 != nil {
		log.Errorf("Get configuration failed with error: %s", err2)

		return err2
	}
	deleteOptions := m.GetMarkDeleteOptions(lctx, rt, rt.ObjectMeta(obj))
	device, err = util.BuildDevice(lctx, conf, device, deleteOptions...)
	if err != nil {
		return err
	}
	d, err := m.UpdateAndReplaceResource(lctx, rt, device.ID, device)
	if err != nil {
		return err
	}
	log.Infof("Moved device: %v", d)
	m.UnsetDeviceInCache(lctx, rt, d)

	return nil
}
