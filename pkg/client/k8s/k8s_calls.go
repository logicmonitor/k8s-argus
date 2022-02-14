package k8s

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// GetAllK8SResources get all k8s resources present in cluster
func GetAllK8SResources(lctx *lmctx.LMContext) (*resourcecache.Store, error) {
	log := lmlog.Logger(lctx)
	tmpStore := resourcecache.NewStore()
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return nil, err
	}
	for _, rt := range enums.ALLResourceTypes {
		if conf.IsMonitoringDisabled(rt) {
			continue
		}

		// check if we have enabled permission for resource type.
		if !permission.HasPermissions(rt) {
			log.Warnf("no permission for resource type %s", rt.String())
			continue
		}

		all, err := GetAndStoreAll(lctx, rt)
		if err != nil {
			return nil, err
		}
		for _, metaObject := range all {
			displayName := util.GetDisplayName(rt, metaObject, conf) //nolint:gosec
			tmpStore.Set(lctx, types.ResourceName{
				Name:     metaObject.Name,
				Resource: rt,
			}, types.ResourceMeta{ // nolint: exhaustivestruct
				Container:   metaObject.Namespace,
				Labels:      metaObject.Labels,
				DisplayName: displayName,
				UID:         metaObject.UID,
			})
		}
	}

	return tmpStore, nil
}

// GetAndStoreAll get
func GetAndStoreAll(lctx *lmctx.LMContext, rt enums.ResourceType) ([]*metav1.PartialObjectMetadata, error) {
	log := lmlog.Logger(lctx)
	result := make([]*metav1.PartialObjectMetadata, 0)
	if rt == enums.ETCD || rt == enums.Unknown {
		return result, nil
	}
	listWatch := cache.NewListWatchFromClient(util.GetK8sRESTClient(config.GetClientSet(), rt.K8SAPIVersion()), rt.String(), corev1.NamespaceAll, fields.Everything())
	listWatch.DisableChunking = true
	list, err := listWatch.List(constants.DefaultListOptions)
	if err != nil {
		log.Warnf("error while extracting resources from cluster:  %v", err)
		return result, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return result, fmt.Errorf("%s: Unable to understand list result %#v (%w)", rt, list, err)
	}
	for _, item := range items {
		accessor, err := meta.Accessor(item)
		if err != nil {
			return nil, err
		}
		result = append(result, meta.AsPartialObjectMetadata(accessor))
		metadata := meta.AsPartialObjectMetadata(accessor)
		now := metav1.Now()
		if metadata.DeletionTimestamp.Before(&now) {
			log.Infof("Ignoring already deleted resource: %v", rt.LMName(metadata))
			continue
		}
		result = append(result, metadata)

	}

	return result, nil

}
