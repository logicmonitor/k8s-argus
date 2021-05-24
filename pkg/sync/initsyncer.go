package sync

import (
	"net/http"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/device"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitSyncer implements the initial sync through logicmonitor API
type InitSyncer struct {
	DeviceManager *device.Manager
}

// Sync sync
func (i *InitSyncer) Sync(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	// Graceful conflicts resolution
	resolveConflicts := false
	clusterGroupID := util.GetClusterGroupID(lctx, i.DeviceManager.LMClient)

	if devicegroup.GetClusterGroupProperty(lctx, constants.ResyncConflictingResourcesProp, i.DeviceManager.LMClient) == "true" {
		resolveConflicts = true
	}
	log.Infof("resolveConflicts is: %v", resolveConflicts)
	defer func() {
		// Reset property so that this would happen gracefully
		devicegroup.DeleteDeviceGroupPropertyByName(lctx, clusterGroupID, &models.EntityProperty{Name: constants.ResyncConflictingResourcesProp, Value: "true"}, i.DeviceManager.LMClient)
	}()
	allK8SResourcesStore := i.DeviceManager.GetAllK8SResources()
	log.Tracef("Resources present on cluster: %v", allK8SResourcesStore)
	resourcesToDelete := map[enums.ResourceType]bool{
		enums.Pods:        true,
		enums.Deployments: true,
		enums.Services:    true,
		enums.Nodes:       true,
		enums.Hpas:        true,
	}

	list := i.DeviceManager.ResourceCache.List()
	log.Tracef("Current cache: %v", list)
	for _, entry := range list {
		log.Tracef("Iterate resource cache entry : %v ", entry)
		cacheResourceName := entry.K
		cacheResourceMeta := entry.V

		childLctx := lmlog.LMContextWithFields(lctx, logrus.Fields{
			"name":  cacheResourceName.Resource.FQName(cacheResourceName.Name),
			"type":  cacheResourceName.Resource.Singular(),
			"ns":    cacheResourceMeta.Container,
			"event": "sync",
		})

		if !resourcesToDelete[cacheResourceName.Resource] {
			continue
		}

		clusterPresentMeta, ok := allK8SResourcesStore.Exists(childLctx, cacheResourceName, cacheResourceMeta.Container)
		if !ok {
			i.deleteDevice(childLctx, log, cacheResourceName, cacheResourceMeta)
		} else {
			i.resolveConflicts(childLctx, resolveConflicts, cacheResourceMeta, clusterPresentMeta, cacheResourceName, log)
		}
	}

	// Flush updated cache to configmaps
	err3 := i.DeviceManager.ResourceCache.Save()
	if err3 != nil {
		log.Errorf("Failed to flush resource cache after resync: %s", err3)
	}
}

// nolint: gocognit
func (i *InitSyncer) resolveConflicts(lctx *lmctx.LMContext, resolveConflicts bool, cacheMeta cache.ResourceMeta, clusterResourceMeta cache.ResourceMeta, cacheResourceName cache.ResourceName, log *logrus.Entry) {
	rt := cacheResourceName.Resource
	if resolveConflicts && clusterResourceMeta.DisplayName != cacheMeta.DisplayName {
		conf, err := config.GetConfig()
		if err != nil {
			log.Errorf("failed to get confing")
			return
		}
		displayNameNew := util.GetDisplayNameNew(rt, &metav1.ObjectMeta{
			Name:      cacheResourceName.Name,
			Namespace: clusterResourceMeta.Container,
		}, conf)
		if cacheMeta.DisplayName != displayNameNew {
			log.Infof("Updating resource by changing displayName to %s", displayNameNew)
			resource, err := i.DeviceManager.FetchDevice(lctx, rt, cacheMeta.LMID)
			if err != nil {
				log.Errorf("failed to fetch resource to change displayname: %v", cacheMeta.LMID)
				return
			}
			options := []types.DeviceOption{
				i.DeviceManager.DisplayName(displayNameNew),
				i.DeviceManager.SystemCategory(rt.GetConflictsCategory(), enums.Delete),
			}
			modifiedResourceValue := *resource
			modifiedResource, err := util.BuildDevice(lctx, conf, &modifiedResourceValue, options...)
			if err != nil {
				log.Errorf("Failed to build modified resource")
				return
			}
			_, err = i.DeviceManager.UpdateAndReplaceResource(lctx, rt, modifiedResource.ID, modifiedResource)
			if err != nil {
				i.handleConflictResource(lctx, err, log, cacheMeta, rt, conf, resource)
				return
			}
			return
		}
		log.Infof("No change in settings to change displayName")
	}
}

func (i *InitSyncer) handleConflictResource(lctx *lmctx.LMContext, err error, log *logrus.Entry, cacheMeta cache.ResourceMeta, rt enums.ResourceType, conf *config.Config, resource *models.Device) {
	deviceDefault := err.(*lm.UpdateDeviceDefault) // nolint: errorlint
	if deviceDefault != nil && deviceDefault.Code() == http.StatusConflict {
		log.Warnf("Still resource conflicts, ignoring resolve conflict for resource")
		if !cacheMeta.HasSysCategory(rt.GetConflictsCategory()) {
			modifiedResource, err := util.BuildDevice(lctx, conf, resource, i.DeviceManager.SystemCategory(rt.GetConflictsCategory(), enums.Add))
			if err != nil {
				log.Errorf("Failed to modify resource to add conflicts category: %s", err)
				return
			}
			_, err = i.DeviceManager.UpdateAndReplaceResource(lctx, rt, modifiedResource.ID, modifiedResource)
			if err != nil {
				log.Errorf("Failed to add conflicts category on resource: %s", err)
			}
		}
		return
	}
	log.Errorf("Failed to modify resource name")
}

func (i *InitSyncer) deleteDevice(lctx *lmctx.LMContext, log *logrus.Entry, resourceName cache.ResourceName, resourceMeta cache.ResourceMeta) {
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}
	if conf.DeleteDevices {
		log.Debugf("Deleting device: %s %v", resourceName.Name, resourceMeta)
		err := i.DeviceManager.DeleteByID(lctx, resourceName.Resource, resourceMeta.LMID)
		if err != nil {
			sc := util.GetHTTPStatusCodeFromLMSDKError(err)
			if sc == http.StatusNotFound {
				log.Tracef("Device does not exist %s, %v", resourceName.Name, resourceMeta.LMID)
				i.DeviceManager.ResourceCache.Unset(resourceName, resourceMeta.Container)
			} else {
				log.Errorf("Failed to delete dangling device %s with ID %v : %s", resourceName.Name, resourceMeta.LMID, err)
			}
		} else {
			i.DeviceManager.ResourceCache.Unset(resourceName, resourceMeta.Container)
			log.Tracef("Deleted dangling device %s with id: %v", resourceName.Name, resourceMeta.LMID)
		}
	} else {
		log.Infof("Soft delete")
		// TODO:: Mark device for deletion if not already
	}
}

// InitSync implements the initial sync through logicmonitor API
func (i *InitSyncer) InitSync(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Start to sync the resource devices")
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}
	clusterName := i.DeviceManager.Base.Config.ClusterName
	// get the cluster info
	parentGroupID := conf.ClusterGroupID
	groupName := util.ClusterGroupName(clusterName)
	rest, err := devicegroup.Find(parentGroupID, groupName, i.DeviceManager.LMClient)
	if err != nil || rest == nil {
		log.Infof("Failed to get the cluster group: %v, parentID: %v", groupName, parentGroupID)
	}
	if rest == nil {
		return
	}

	log.Infof("Finished syncing the resource devices")
}

// RunPeriodicSync runs synchronization periodically.
func (i *InitSyncer) RunPeriodicSync(syncTime time.Duration) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "periodic-sync"}))
	go func() {
		for {
			time.Sleep(syncTime)
			i.Sync(lctx)
		}
	}()
}
