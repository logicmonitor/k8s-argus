package namespace

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

var rt = enums.Namespaces

// OldWatcher represents a watcher type that watches namespaces.
type OldWatcher struct {
	ResourceGroups map[string]int32
	mu             sync.RWMutex
	types.ResourceManager
	ResourceCache types.ResourceCache
	*types.LMRequester
}

func NewOldWatcher(manager types.ResourceManager, resourceCache types.ResourceCache, lmRequester *types.LMRequester) *OldWatcher {
	watcher := &OldWatcher{
		ResourceGroups:  make(map[string]int32),
		ResourceManager: manager,
		ResourceCache:   resourceCache,
		LMRequester:     lmRequester,
		mu:              sync.RWMutex{},
	}
	resourceCache.AddCacheHook(types.CacheHook{
		Hook:      getHook(watcher, types.CacheSet),
		Predicate: getHookPredicate(types.CacheSet),
	})
	resourceCache.AddCacheHook(types.CacheHook{
		Hook:      getHook(watcher, types.CacheUnset),
		Predicate: getHookPredicate(types.CacheUnset),
	})
	return watcher
}

func getHook(watcher *OldWatcher, action types.CacheAction) func(rn types.ResourceName, meta types.ResourceMeta) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"cache_hook": "namespace", "action": action.String()}))
	log := lmlog.Logger(lctx)
	return func(rn types.ResourceName, meta types.ResourceMeta) {
		watcher.mu.Lock()
		defer watcher.mu.Unlock()
		log.Tracef("Hook %s called for: %s %d", action, rn.Name, meta.LMID)
		if action == types.CacheSet {
			watcher.ResourceGroups[rn.Name] = meta.LMID
		} else if action == types.CacheUnset {
			delete(watcher.ResourceGroups, rn.Name)
		}
	}
}

func getHookPredicate(expectedAction types.CacheAction) func(action types.CacheAction, rn types.ResourceName, meta types.ResourceMeta) bool {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"cache_hook_predicate": "namespace", "action": expectedAction.String()}))
	log := lmlog.Logger(lctx)
	return func(action types.CacheAction, rn types.ResourceName, meta types.ResourceMeta) bool {
		ok := false
		if action == expectedAction && rn.Resource == enums.Namespaces {
			for _, e := range enums.ALLResourceTypes {
				if e.IsNamespaceScopedResource() && e.TitlePlural() == rn.Name {
					ok = true
					break
				}
			}
		}
		log.Tracef("Evaluating %s hook predicate for %s %s: %v", expectedAction, rn.Resource, rn.Name, ok)

		return ok
	}
}

// ResourceType resource
func (w *OldWatcher) ResourceType() enums.ResourceType {
	return enums.Namespaces
}

// GetConfig get
func (w *OldWatcher) GetConfig() *types.WConfig {
	return nil
}

// AddFunc is a function that implements the Watcher interface.
func (w *OldWatcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": namespace.Name, "event": "add", "type": rt.Singular()}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling add namespace event: %s", namespace.Name)

		w.createNamspaceResourceGroupTree(lctx, namespace)
	}
}

func (w *OldWatcher) createNamspaceResourceGroupTree(lctx *lmctx.LMContext, namespace *corev1.Namespace) {
	conf, err := config.GetConfig()
	if err != nil {
		return
	}
	if conf.EnableNewResourceTree {
		// this will be based on new resource tree where namespace groups will be created and all resources to put under it
		w.createNewResourceGroupTree(lctx, namespace)
	} else {
		// resource wise separate static groups and underneath namespace groups in each
		w.createPreviousResourceGroupTree(lctx, namespace)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *OldWatcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		namespace := newObj.(*corev1.Namespace) // nolint: forcetypeassert
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": namespace.Name, "event": "update", "type": rt.Singular()}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling update namespace event: %s", namespace.Name)
		w.createNamspaceResourceGroupTree(lctx, namespace)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *OldWatcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		conf, err := config.GetConfig()
		if err != nil {
			return
		}
		namespace := obj.(*corev1.Namespace) // nolint: forcetypeassert
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"name": namespace.Name, "event": "delete", "type": rt.Singular()}))
		lctx = lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})

		log := lmlog.Logger(lctx)
		log.Debugf("Handle delete namespace event: %s", namespace.Name)

		metaList, ok := w.ResourceCache.Get(lctx, types.ResourceName{
			Name:     namespace.Name,
			Resource: enums.Namespaces,
		})
		if !ok {
			log.Errorf("No resource groups exists in cache for namespace : %s", namespace.Name)
			return
		}
		for _, meta := range metaList {
			clctx := lmlog.LMContextWithLMResourceID(lctx, meta.LMID)
			err := w.DeleteResourceGroup(clctx, enums.Namespaces, meta.LMID, false)
			if err != nil {
				log.Errorf("Failed to delete resource group for namespace [%d] of parent [%s]: %s", meta.LMID, meta.Container, err)
			} else {
				log.Infof("Deleted resource group [%d] of parent [%s]", meta.LMID, meta.Container)
			}
		}
	}
}

// nolint: unused
func (w *OldWatcher) createNewResourceGroupTree(lctx *lmctx.LMContext, namespace *corev1.Namespace) {
	log := lmlog.Logger(lctx)
	meta, ok := w.ResourceCache.Get(lctx, types.ResourceName{
		Name:     "Namespaces",
		Resource: enums.Namespaces,
	})
	if !ok {
		log.Errorf("cannot find \"Namespaces\" resource group to add \"%s\"", namespace.Name)
		return
	}

	conf, err := config.GetConfig()
	if err != nil {
		return
	}

	deletedBuilder := dgbuilder.NewAppliesToBuilder().
		Auto("namespace").Equals(namespace.Name).And().
		Auto("clustername").Equals(conf.ClusterName).And().
		OpenBracket()
	for _, e := range enums.ALLResourceTypes {
		if e == enums.Namespaces || !e.IsNamespaceScopedResource() {
			continue
		}
		deletedBuilder = deletedBuilder.HasCategory(e.GetDeletedCategory()).Or()
	}
	deletedBuilder.TrimOrCloseBracket()

	log.Debugf("deleted applies to: %s", deletedBuilder.Build())

	resourceTree := &types.ResourceGroupTree{
		Options: []types.ResourceGroupOption{
			w.ResourceManager.GroupName(namespace.Name),
			w.ResourceManager.ParentID(meta[0].LMID),
			w.ResourceManager.DisableAlerting(conf.DisableAlerting),
			w.ResourceManager.AppliesTo(dgbuilder.NewAppliesToBuilder().
				Auto("namespace").Equals(namespace.Name).And().
				Auto("clustername").Equals(conf.ClusterName),
			),
		},
		ChildGroups: []*types.ResourceGroupTree{
			{
				// Set operation = (A'B)' ## if hard delete resources is false and enable ns specific _deleted group then on create else no
				DontCreate: !(!conf.DeleteResources && conf.EnableNamespacesDeletedGroups),
				Options: []types.ResourceGroupOption{
					w.ResourceManager.GroupName(constants.DeletedResourceGroup),
					w.ResourceManager.DisableAlerting(true),
					w.ResourceManager.AppliesTo(deletedBuilder),
				},
			},
		},
	}
	err = w.CreateResourceGroupTree(lctx, resourceTree, false)
	if err != nil {
		return
	}
}

func (w *OldWatcher) createPreviousResourceGroupTree(lctx *lmctx.LMContext, namespace *corev1.Namespace) {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig()
	if err != nil {
		log.Errorf("Failed to get config")
		return
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	log.Debugf("Creating resource groups for namespace [%s] under parent resource groups: %v", namespace.Name, w.ResourceGroups)
	for resourceGroupName, parentID := range w.ResourceGroups {
		rt, err := enums.ParseResourceType(resourceGroupName)
		if err != nil || rt == enums.Unknown {
			continue
		}
		if !rt.IsNamespaceScopedResource() {
			continue
		}

		options := &types.ResourceGroupTree{
			Options: []types.ResourceGroupOption{
				w.GroupName(namespace.Name),
				w.ParentID(parentID),
				w.AppliesTo(dgbuilder.NewAppliesToBuilder().HasCategory(rt.GetCategory()).And().Auto("namespace").Equals(namespace.Name).And().Auto("clustername").Equals(conf.ClusterName)),
				w.DisableAlerting(conf.DisableAlerting),
			},
		}
		err = w.CreateResourceGroupTree(lctx, options, false)
		if err != nil {
			log.Errorf("Failed to add %q namespace group under %q resource group. Error: %v", namespace.Name, resourceGroupName, err)
			continue
		}
	}
	log.Debugf("Created resource groups for namespace: %s", namespace.Name)
}
