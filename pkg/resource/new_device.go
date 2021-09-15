package resource

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// Add implements types.ResourceManager.
func (m *Manager) Add(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.ResourceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return nil, err
	}
	resource, err := util.BuildResource(lctx, conf, nil, options...)
	if err != nil {
		return nil, fmt.Errorf("could not build resource object: %w", err)
	}

	objectMeta, _ := rt.ObjectMeta(obj)
	warn, valid, err2 := validateNewResource(lctx, rt, resource, objectMeta)
	if !valid {
		return nil, err2
	}
	// if warning, then just logging and proceeding further
	if warn {
		log.Warnf("resource validation error %s", err2)
	}

	log.Debugf("Does resource exists in cache")

	_, ok := m.DoesResourceExistInCache(lctx, rt, resource, false)
	if ok {
		return nil, fmt.Errorf("resource already present, ignoring add event: %w", aerrors.ErrResourceExists)
	}

	resource, err = util.BuildResource(lctx, conf, nil, options...)
	if err != nil {
		return nil, fmt.Errorf("could not build resource object: %w", err)
	}

	resultResource, err3 := m.addResource(lctx, rt, resource)
	if err3 != nil {
		log.Warnf("add resource failed with: %s", err3.Error())
		m2, err4 := m.handleConflict(lctx, rt, obj, err3, resource, options...)
		if err4 != nil {
			return m2, err4
		}
	}

	return resultResource, nil
}

func (m *Manager) handleConflict(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, err3 error, resource *models.Device, options ...types.ResourceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	var resourceDefault *lm.AddDeviceDefault
	ok := errors.Is(err3, resourceDefault)
	if !ok {
		return nil, err3
	}
	// Conflicts due to cache miss or outside cluster same resource
	// 99.99% call should not reach here: more the control here, more the portal requests
	if resourceDefault != nil && resourceDefault.Code() == http.StatusConflict {
		// nolint: godox
		existingResource, err := m.FindByDisplayName(lctx, rt, *resource.DisplayName)
		if err != nil {
			// return error earlier received
			return nil, fmt.Errorf("failed with error: %w. Also cannot find resource with display name: %s", err3, *resource.DisplayName)
		}
		if util.GetResourcePropertyValue(existingResource, constants.K8sResourceUIDPropertyKey) != "" &&
			util.GetResourcePropertyValue(existingResource, constants.K8sResourceUIDPropertyKey) == util.GetResourcePropertyValue(resource, constants.K8sResourceUIDPropertyKey) {
			// TODO: PROM_METRIC cache miss metrics, fetch and correlate resource - as full name is enforced now
			conf, err5 := config.GetConfig(lctx)
			if err5 != nil {
				return nil, err5
			}
			resource, err6 := util.BuildResource(lctx, conf, existingResource, options...)
			if err6 != nil {
				log.Errorf("Failed to buid resource object")
				return nil, err6
			}
			resultResource2, err4 := m.UpdateAndReplaceResource(lctx, rt, existingResource.ID, resource)
			if err4 != nil {
				log.Errorf("failed to add update resource: %s", err4)
				// nolint: godox
				// TODO: PROM_METRIC Add resource failed metric
				return nil, err4
			}

			return resultResource2, nil
		}
		log.Errorf("Conflicts with other resource that is not part of this cluster")
		// nolint: godox
		// TODO: PROM_METRIC Add resource failed metric

		return nil, err3
	}

	return nil, err3
}

// Update update
// nolint: cyclop
func (m *Manager) Update(lctx *lmctx.LMContext, rt enums.ResourceType, oldObj interface{}, newObj interface{}, options ...types.ResourceOption) (*models.Device, error) {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return nil, err
	}
	resource, err := util.BuildResource(lctx, conf, nil, options...)
	if err != nil {
		return nil, err
	}
	ce, ok := m.DoesResourceExistInCache(lctx, rt, resource, false)
	if !ok {
		log.Debugf("Resource does not exist in cache")

		return nil, aerrors.ErrCacheMiss
	}
	// If the LM resource is modified outside of argus:
	// For ex:
	// 1. User manually added resource property on LM resource 2. PropertySources added property on LM resource
	// Partial patch API is not allowed - Either it will replace all set of customProperties or nothing
	// Similarly, system.categories is comma separated string, so to add or modify, we need to fetch current string from LM resource

	fetchedResource, err := m.FetchResource(lctx, rt, ce.LMID)
	if err != nil {
		// This indicates that the lm resource id present in cache is wrong, consider as cache miss
		// nolint: godox
		// TODO: PROM_METRIC cache miss gauge
		log.Errorf("Failed to fetch existing resource: %s", err)

		return nil, err
	}

	// If UID mismatch then fail with error, auto prop is not going to change, and even if it updates then metrics of prev and new get mixed up
	if util.GetResourcePropertyValue(fetchedResource, constants.K8sResourceUIDPropertyKey) != util.GetResourcePropertyValue(resource, constants.K8sResourceUIDPropertyKey) {
		return nil, fmt.Errorf("mismatch in UID: previous (%s), new (%s)", util.GetResourcePropertyValue(fetchedResource, constants.K8sResourceUIDPropertyKey), util.GetResourcePropertyValue(resource, constants.K8sResourceUIDPropertyKey))
	}

	// Put ID for update request
	fetchedResource.ID = ce.LMID
	// retain displayName
	options = append(options, m.DisplayName(*fetchedResource.DisplayName))
	prevIP := util.GetResourcePropertyValue(fetchedResource, constants.K8sSystemIPsPropertyKey)
	fetchedResource, err = util.BuildResource(lctx, conf, fetchedResource, options...)
	if err != nil {
		return nil, err
	}
	newIP := util.GetResourcePropertyValue(fetchedResource, constants.K8sSystemIPsPropertyKey)
	if rt.IsK8SPingResource() && prevIP != newIP {
		log.Infof("ResourceIP [%s] has changed, updating new ip [%s]", prevIP, newIP)
		// set name to ip
		fetchedResource, err = util.BuildResource(lctx, conf, fetchedResource, m.Name(newIP))
		if err != nil {
			return nil, fmt.Errorf("failed to build resource to update New IP: %w", err)
		}

		// patch with ip as hostname
		fetchedResource, err = m.PatchResource(lctx, rt, fetchedResource, constants.NameFieldName)
		if err != nil {
			return nil, fmt.Errorf("failed to update resource with New IP: %w", err)
		}

		// wait to copy hostname into system.ips
		ctx, cancel := context.WithTimeout(context.Background(), *conf.SysIpsWaitTimeout)
		defer cancel()
		fetchedResource, err = m.WaitToReflectSysIps(ctx, lctx, rt, fetchedResource, newIP)
		if err != nil {
			return nil, fmt.Errorf("failed while waiting to update resource with New IP: %w", err)
		}
		fetchedResource, err = util.BuildResource(lctx, conf, fetchedResource, options...)
		if err != nil {
			return nil, fmt.Errorf("failed to build resource after updating new IP: %w", err)
		}
	}
	// Update the resource
	return m.UpdateAndReplaceResource(lctx, rt, fetchedResource.ID, fetchedResource)
}

// Delete delete
func (m *Manager) Delete(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.ResourceOption) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return err
	}
	resource, err := util.BuildResource(lctx, conf, nil, options...)
	if err != nil {
		return err
	}
	ce, ok := m.DoesResourceExistInCache(lctx, rt, resource, false)
	delLctx := lmlog.LMContextWithLMResourceID(lctx, ce.LMID)
	log = lmlog.Logger(delLctx)

	if !ok {
		log.Tracef("Resource does not exist in cache to delete")
		return nil
	}
	if ce.LMID == 0 {
		m.UnsetResourceInCache(delLctx, rt, resource)

		return nil
	}
	resource.ID = ce.LMID
	err = m.deleteResource(delLctx, rt, resource)
	if err != nil {
		return err
	}

	return nil
}

// MarkDeleted mark
func (m *Manager) MarkDeleted(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}, options ...types.ResourceOption) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return err
	}
	resource, err := util.BuildResource(lctx, conf, nil, options...)
	if err != nil {
		return err
	}
	ce, ok := m.DoesResourceExistInCache(lctx, rt, resource, false)
	if !ok {
		return nil
	}
	if ce.LMID == 0 {
		m.UnsetResourceInCache(lctx, rt, resource)

		return nil
	}
	delLctx := lmlog.LMContextWithLMResourceID(lctx, ce.LMID)
	log = lmlog.Logger(delLctx)
	resource.ID = ce.LMID
	conf, err2 := config.GetConfig(lctx)
	if err2 != nil {
		log.Errorf("Get configuration failed with error: %s", err2)

		return err2
	}
	meta, _ := rt.ObjectMeta(obj)
	deleteOptions := m.GetMarkDeleteOptions(delLctx, rt, meta)
	resource, err = util.BuildResource(delLctx, conf, resource, deleteOptions...)
	if err != nil {
		return err
	}
	d, err := m.UpdateAndReplaceResource(delLctx, rt, resource.ID, resource)
	if err != nil {
		return err
	}

	m.UnsetResourceInCache(delLctx, rt, d)

	return nil
}

func (m *Manager) WaitToReflectSysIps(ctx context.Context, lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device, ip string) (*models.Device, error) {
	backOff := 5 * time.Second // nolint: gomnd
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting to change system.ips property on resource [%s]: %w", *resource.DisplayName, ctx.Err())

		default:
			isIPUpdated, fetchResource, err := m.doesResourceHaveExpectedIP(lctx, resource, rt, ip)
			if err != nil && errors.Is(err, aerrors.ErrInvalidCache) {
				return fetchResource, fmt.Errorf("resource [%s] got deleted in between while argus waiting for system.ips to reflect: %w", *resource.DisplayName, err)
			}
			if isIPUpdated {
				return fetchResource, nil
			}
			time.Sleep(backOff)
			backOff *= 2
		}
	}
}

func (m *Manager) doesResourceHaveExpectedIP(lctx *lmctx.LMContext, resource *models.Device, rt enums.ResourceType, expectedPodIP string) (bool, *models.Device, error) {
	params := lm.NewGetDeviceByIDParams()
	params.SetID(resource.ID)
	command := m.GetResourceByIDCommand(lctx, params)
	resp, err := m.SendReceive(lctx, command)
	if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		m.UnsetResourceInCache(lctx, rt, resource)
		return false, nil, fmt.Errorf("resource [%s] does not exist: %w", *resource.DisplayName, aerrors.ErrInvalidCache)
	}
	if err != nil {
		return false, nil, err
	}

	fetchedResource := resp.(*lm.GetDeviceByIDOK).Payload
	podIP := util.GetResourcePropertyValue(fetchedResource, constants.K8sSystemIPsPropertyKey)
	return podIP == expectedPodIP, fetchedResource, err
}
