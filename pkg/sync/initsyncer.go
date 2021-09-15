package sync

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/client/k8s"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitSyncer implements the initial sync through logicmonitor API
type InitSyncer struct {
	*types.LMRequester
	ResourceManager types.ResourceManager
}

// Sync sync
// nolint: cyclop
func (i *InitSyncer) Sync(lctx *lmctx.LMContext) {
	defer metrics.ObserveTime(metrics.StartTimeObserver(metrics.SyncTimeSummary))
	log := lmlog.Logger(lctx)
	// Graceful conflicts resolution
	resolveConflicts := false
	clusterGroupID, err := util.GetClusterGroupID(lctx, i.LMRequester)
	if err != nil {
		log.Error(err.Error())

		return
	}

	if resourcegroup.GetClusterGroupProperty(lctx, constants.ResyncConflictingResourcesProp, i.LMRequester) == "true" {
		resolveConflicts = true
	}
	log.Infof("resolveConflicts is: %v", resolveConflicts)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}
	defer func() {
		childLctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
		// Reset property so that this would happen gracefully
		resourcegroup.DeleteResourceGroupPropertyByName(childLctx, clusterGroupID, &models.EntityProperty{Name: constants.ResyncConflictingResourcesProp, Value: "true"}, i.LMRequester)
	}()
	allK8SResourcesStore, err := k8s.GetAllK8SResources(lctx)
	if err != nil {
		log.Errorf("Failed to fetch current resource present on cluster: %s", err)
		return
	}
	log.Tracef("Resources present on cluster: %v", allK8SResourcesStore)
	ignoreSync := map[enums.ResourceType]bool{
		enums.ETCD:    true,
		enums.Unknown: true,
	}

	list := i.ResourceManager.GetResourceCache().List()
	log.Tracef("Current cache: %v", list)

	for _, entry := range list {
		log.Tracef("Iterate resource cache entry : %v ", entry)
		cacheResourceName := entry.K
		cacheResourceMeta := entry.V

		if ignoreSync[cacheResourceName.Resource] {
			continue
		}
		childLctx := lmlog.LMContextWithFields(lctx, logrus.Fields{
			"name":  cacheResourceName.Resource.FQName(cacheResourceName.Name),
			"type":  cacheResourceName.Resource.Singular(),
			"ns":    cacheResourceMeta.Container,
			"event": "sync",
		})
		childLctx = childLctx.LMContextWith(map[string]interface{}{constants.PartitionKey: fmt.Sprintf("%s-%s", cacheResourceName.Resource.String(), cacheResourceName.Name)})

		if cacheResourceName.Resource == enums.Namespaces {
			if err := i.deleteNamespace(allK8SResourcesStore, childLctx, cacheResourceName, cacheResourceMeta, log, conf); err != nil && !errors.Is(err, aerrors.ErrResourceGroupIsNotEmpty) &&
				!errors.Is(err, aerrors.ErrResourceGroupParentIsNotValid) &&
				!strings.Contains(err.Error(), util.ClusterGroupName(conf.ClusterName)) {
				log.Errorf("failed to delete resource group: %s", err)
			}
			continue
		}
		clusterPresentMeta, ok := allK8SResourcesStore.Exists(childLctx, cacheResourceName, cacheResourceMeta.Container)
		// Delete resource if no more exists or delete if UID does not match.
		if !ok ||
			clusterPresentMeta.UID != cacheResourceMeta.UID ||
			(conf.RegisterGenericFilter && !util.EvaluateExclusion(clusterPresentMeta.Labels)) {

			i.deleteResource(childLctx, log, cacheResourceName, cacheResourceMeta)
		} else if resolveConflicts {
			i.resolveConflicts(childLctx, cacheResourceMeta, clusterPresentMeta, cacheResourceName, log)
		}
	}

	// Flush updated cache to configmaps
	err3 := i.ResourceManager.GetResourceCache().Save(lctx)
	if err3 != nil {
		log.Errorf("Failed to flush resource cache after resync: %s", err3)
	}
}

func (i *InitSyncer) deleteNamespace(allK8SResourcesStore *resourcecache.Store, childLctx *lmctx.LMContext, cacheResourceName types.ResourceName, cacheResourceMeta types.ResourceMeta, log *logrus.Entry, conf *config.Config) error {
	if _, ok := allK8SResourcesStore.Get(childLctx, cacheResourceName); !ok {
		list := i.ResourceManager.GetResourceCache().ListWithFilter(func(k types.ResourceName, v types.ResourceMeta) bool {
			return fmt.Sprintf("%d", v.LMID) == cacheResourceMeta.Container
		})
		if len(list) == 0 {
			return fmt.Errorf("failed to determine parent group of resource group %s [%d]", cacheResourceName.Name, cacheResourceMeta.LMID)
		}
		if len(list) > 1 {
			return fmt.Errorf("more than one parent group of resource group %s [%d]", cacheResourceName.Name, cacheResourceMeta.LMID)
		}
		e, err := enums.ParseResourceType(list[0].K.Name)
		if err != nil {
			return fmt.Errorf("%w", aerrors.ErrResourceGroupParentIsNotValid)
		}
		if (conf.EnableNewResourceTree && e == enums.Namespaces) ||
			(!conf.EnableNewResourceTree && e != enums.Namespaces && e.IsNamespaceScopedResource()) {
			return i.deleteResourceGroup(childLctx, cacheResourceName, cacheResourceMeta)
		}
	}
	return nil
}

// nolint: gocognit
func (i *InitSyncer) resolveConflicts(lctx *lmctx.LMContext, cacheMeta types.ResourceMeta, clusterResourceMeta types.ResourceMeta, cacheResourceName types.ResourceName, log *logrus.Entry) {
	rt := cacheResourceName.Resource
	if clusterResourceMeta.DisplayName != cacheMeta.DisplayName || cacheMeta.HasSysCategory(rt.GetConflictsCategory()) {
		conf, err := config.GetConfig(lctx)
		if err != nil {
			log.Errorf("failed to get confing")
			return
		}
		displayNameNew := util.GetDisplayName(rt, meta.AsPartialObjectMetadata(&metav1.ObjectMeta{
			Name:      cacheResourceName.Name,
			Namespace: clusterResourceMeta.Container,
		}), conf)
		if cacheMeta.DisplayName != displayNameNew || cacheMeta.HasSysCategory(rt.GetConflictsCategory()) {
			log.Infof("Updating resource by changing displayName to %s", displayNameNew)
			options := []types.ResourceOption{
				i.ResourceManager.DisplayName(displayNameNew),
				i.ResourceManager.SystemCategory(rt.GetConflictsCategory(), enums.Delete),
			}
			_, err = i.ResourceManager.UpdateResourceByID(lctx, rt, cacheMeta.LMID, options...)
			if err != nil {
				log.Errorf("Failed to update resource with error: %s", err)

				return
			}

			return
		}
		log.Infof("No change in settings to change displayName")
	}
}

func (i *InitSyncer) deleteResource(lctx *lmctx.LMContext, log *logrus.Entry, resourceName types.ResourceName, resourceMeta types.ResourceMeta) {
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}
	if conf.DeleteResources &&
		!util.IsArgusPodCacheMeta(lctx, resourceName.Resource, resourceMeta) {
		log.Info("Deleting resource")
		err := i.ResourceManager.DeleteResourceByID(lctx, resourceName.Resource, resourceMeta.LMID)
		if err != nil {
			sc := util.GetHTTPStatusCodeFromLMSDKError(err)
			if sc == http.StatusNotFound {
				log.Tracef("Resource does not exist %s, %v", resourceName.Name, resourceMeta.LMID)
				i.ResourceManager.GetResourceCache().Unset(lctx, resourceName, resourceMeta.Container)
			} else {
				log.Errorf("Failed to delete dangling resource %s with ID %v : %s", resourceName.Name, resourceMeta.LMID, err)
			}
		} else {
			i.ResourceManager.GetResourceCache().Unset(lctx, resourceName, resourceMeta.Container)
			log.Tracef("Deleted dangling resource %s with id: %v", resourceName.Name, resourceMeta.LMID)
		}
	} else {
		log.Info("Soft delete")
		deleteOptions := i.ResourceManager.GetMarkDeleteOptions(lctx, resourceName.Resource, meta.AsPartialObjectMetadata(&metav1.ObjectMeta{}))

		_, err = i.ResourceManager.UpdateResourceByID(lctx, resourceName.Resource, resourceMeta.LMID, deleteOptions...)
		if err != nil {
			log.Errorf("failed to mark resource as deleted: %s", err)
		} else {
			log.Infof("Marked resource as deleted")
		}
	}
}

// RunPeriodicSync runs synchronization periodically.
func (i *InitSyncer) RunPeriodicSync() {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": "periodic-sync"}))

	go func() {
		for {
			conf, err := config.GetConfig(lctx)
			if err != nil {
				time.Sleep(constants.DefaultPeriodicDeleteInterval) // nolint: gomnd
			} else {
				time.Sleep(*conf.Intervals.PeriodicDeleteInterval)
			}
			i.Sync(lctx)
		}
	}()
}

func (i *InitSyncer) deleteResourceGroup(lctx *lmctx.LMContext, name types.ResourceName, resourceMeta types.ResourceMeta) error {
	err := i.ResourceManager.DeleteResourceGroup(lctx, name.Resource, resourceMeta.LMID, true)
	if err != nil && !errors.Is(err, aerrors.ErrResourceGroupIsNotEmpty) {
		return fmt.Errorf("failed to delete resource group: %w", err)
	}
	if err != nil && errors.Is(err, aerrors.ErrResourceGroupIsNotEmpty) {
		return err
	}
	return nil
}
