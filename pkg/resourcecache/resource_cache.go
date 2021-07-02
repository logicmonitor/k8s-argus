package resourcecache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/promperf"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	// Version to differentiate cache flushed to configmap is with previous datastructures or not.
	// Whenever cache.ResourceName and cache.ResourceMeta struct modifies, increase the version so that it will easier to parse previously created CMs and convert it to newer data structure
	Version = "v1"

	// CMMaxBytes constant to set max threshold for byte size of configmap
	// 1MiB is the capacity of configmap data, but we are storing 80% of the data in a chunk
	CMMaxBytes = 1048576 * 80 / 100
)

const rateLimitBackoffTime = 10 * time.Second

// ResourceCache to maintain a resource cache to calculate delta between resource presence on server and on cluster
type ResourceCache struct {
	store         *Store
	rwm           sync.RWMutex
	rebuildMutex  sync.Mutex
	resyncPeriod  time.Duration
	flushTimeToCM time.Duration
	flushMU       sync.Mutex
	*types.LMRequester
	clusterGrpID int32
	stateLoaded  bool

	// On restart, first dump loaded from config map
	dumpID int64

	// soft refresh last time
	softRefreshLast map[string]time.Time
	softRefreshMu   sync.Mutex

	hooks   []types.CacheHook
	hookrwm sync.RWMutex
}

// NewResourceCache create new ResourceCache object
func NewResourceCache(facadeObj *types.LMRequester, rp time.Duration) *ResourceCache {
	resourceCache := &ResourceCache{
		store:         NewStore(),
		rwm:           sync.RWMutex{},
		resyncPeriod:  rp,
		flushTimeToCM: 1 * time.Minute,
		flushMU:       sync.Mutex{},
		LMRequester:   facadeObj,
		clusterGrpID:  -1,
		stateLoaded:   false,
		rebuildMutex:  sync.Mutex{},
	}

	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"resource_cache": "init"}))
	log := lmlog.Logger(lctx)
	if err := resourceCache.Load(lctx); err != nil {
		log.Warnf("ResourceCache is not loaded from its last state stored in configmaps, it may take more time to rebuild cache on first run. Error: %s", err)
	} else {
		resourceCache.stateLoaded = true
		log.Info("ResourceCache loaded from configmaps")
	}

	return resourceCache
}

func (rc *ResourceCache) AddCacheHook(hook types.CacheHook) {
	rc.hookrwm.Lock()
	defer rc.hookrwm.Unlock()
	rc.hooks = append(rc.hooks, hook)
	// run hook on existing items
	for _, item := range rc.List() {
		if hook.Predicate(types.CacheSet, item.K, item.V) {
			hook.Hook(item.K, item.V)
		}
	}
}

// Run start a goroutine to resync cache periodically with server
func (rc *ResourceCache) Run() {
	initialised := make(chan bool)
	go rc.AutoCacheBuilder(initialised)

	// Periodically flush cache to configmap for recovery after application restart or pod restart
	go rc.CacheToConfigMapDumper()

	// Wait for first initialization
	<-initialised
}

func (rc *ResourceCache) AutoCacheBuilder(initialised chan<- bool) {
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "resource",
			Subsystem:   "cache",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "auto"},
		})

	if !rc.stateLoaded {
		rc.rebuildCache(gauge)
		initialised <- true
	} else {
		initialised <- true
	}
	close(initialised)

	for {
		// to keep constant interval between cache rebuild runs, as rebuild cache is heavy operation so ticker may lead back to back runs
		time.Sleep(rc.resyncPeriod)
		rc.rebuildCache(gauge)
	}
}

func (rc *ResourceCache) CacheToConfigMapDumper() {
	func(resourceCache *ResourceCache) {
		for range time.Tick(resourceCache.flushTimeToCM) {
			debugID := util.GetShortUUID()
			lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"debug_id": debugID}))
			log := lmlog.Logger(lctx)
			start := time.Now()
			err := resourceCache.Save()
			if err != nil {
				log.Errorf("Flush cache to cm failed: %s", err)
			} else {
				log.Infof("Cache flushed to configmaps in time %v", time.Since(start))
			}
		}
	}(rc)
}

// Rebuild graceful rebuild cache
func (rc *ResourceCache) Rebuild(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Gracefully building cache")
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "cache",
			Subsystem:   "resource",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "graceful"},
		})
	rc.rebuildCache(gauge)
}

func (rc *ResourceCache) rebuildCache(gauge prometheus.Gauge) {
	rc.rebuildMutex.Lock()
	defer rc.rebuildMutex.Unlock()
	defer promperf.Duration(promperf.Track(gauge))
	debugID := util.GetShortUUID()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"debug_id": debugID}))
	log := lmlog.Logger(lctx)
	log.Infof("Resource cache fetching resources")
	resources, err := rc.getAllResources(lctx)
	if resources == nil || err != nil {
		log.Errorf("Failed to fetch resources")
	} else {
		log.Debugf("Resync cache map")
		rc.resetCacheStore(resources)
		log.Debugf("Resync cache done")
	}

	if err := rc.Save(); err != nil {
		log.Errorf("cache to cm failed: %s", err)
	} else {
		log.Infof("Cache size: %v", rc.store.Size())
	}
}

// DeviceGroupData schema to send resource groups from api response
type DeviceGroupData struct {
	ResourceName  string
	NamespaceName int32
	LMID          int32
	CustomProps   []*models.NameAndValue
}

func (rc *ResourceCache) getAllResources(lctx *lmctx.LMContext) (*Store, error) {
	log := lmlog.Logger(lctx)
	clusterGroupID := util.GetClusterGroupID(lctx, rc.LMRequester)

	if clusterGroupID == -1 {
		err := fmt.Errorf("no cluster resource group found")
		log.Errorf(err.Error())

		return nil, err
	}
	tmpStore := NewStore()
	resourceChan := make(chan *models.Device)
	resourceGroupChan := make(chan DeviceGroupData)
	resourceFinished := make(chan bool)
	resourceGroupFinished := make(chan bool)

	go rc.accumulateDeviceCache(lctx, resourceChan, tmpStore, resourceFinished)
	go rc.accumulateDeviceGroupCache(lctx, resourceGroupChan, tmpStore, resourceGroupFinished)

	grpIDChan := make(chan int32)

	go rc.fetchGroupDevices(lctx, rc.LMRequester, grpIDChan, resourceChan)

	grpIDChan <- clusterGroupID
	if conf, err := config.GetConfig(); err == nil {
		g, err := rc.getDeviceGroupByID(lctx, clusterGroupID)
		if err != nil {
			return nil, err
		}
		resourceGroupChan <- DeviceGroupData{
			ResourceName:  util.ClusterGroupName(conf.ClusterName),
			NamespaceName: conf.ClusterGroupID,
			LMID:          clusterGroupID,
			CustomProps:   g.Payload.CustomProperties,
		}
	}
	start := time.Now()
	conf, err := config.GetConfig()
	if err == nil && conf.ResourceContainerGroupID != nil {
		grpIDChan <- *conf.ResourceContainerGroupID
	}
	rc.getAllGroups(lctx, clusterGroupID, grpIDChan, resourceGroupChan, 2)
	log.Infof("Resource group fetched in %v seconds", time.Since(start).Seconds())
	close(grpIDChan)
	close(resourceGroupChan)
	<-resourceGroupFinished
	<-resourceFinished

	return tmpStore, nil
}

func (rc *ResourceCache) getDevices(lctx *lmctx.LMContext, grpID int32) (*lm.GetImmediateDeviceListByDeviceGroupIDOK, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
	params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
	params.SetID(grpID)
	command := rc.LMRequester.GetImmediateResourceListByResourceGroupIDCommand(clctx, params)
	resp, err := rc.LMRequester.SendReceive(clctx, command)
	if err != nil {
		return nil, err
	}

	return resp.(*lm.GetImmediateDeviceListByDeviceGroupIDOK), nil
}

func (rc *ResourceCache) fetchGroupDevices(lctx *lmctx.LMContext, b *types.LMRequester, inChan <-chan int32, outChan chan<- *models.Device) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		close(outChan)
		log.Infof("Resource fetch completed in %v seconds", time.Since(start).Seconds())
	}()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go rc.ResourceGroupProcessor(lctx, inChan, outChan, &wg)
	}
	wg.Wait()
}

func (rc *ResourceCache) ResourceGroupProcessor(lctx *lmctx.LMContext, inChan <-chan int32, outChan chan<- *models.Device, wg *sync.WaitGroup) {
	log := lmlog.Logger(lctx)
	defer func() {
		wg.Done()
	}()
	for grpID := range inChan {
		log.Debugf("Fetching resources of group %v", grpID)

		resp, err := rc.getDevices(lctx, grpID)
		if err != nil {
			log.Warnf("fetch resources for %v failed with %v", grpID, err)
		}
		if resp != nil && util.GetHTTPStatusCodeFromLMSDKError(resp) == http.StatusOK {
			for _, resource := range resp.Payload.Items {
				outChan <- resource
			}
		}
	}
}

func (rc *ResourceCache) accumulateDeviceCache(lctx *lmctx.LMContext, inChan <-chan *models.Device, store *Store, finished chan<- bool) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		finished <- true
		close(finished)
		log.Infof("Resource accumulation completed in %v seconds", time.Since(start).Seconds())
	}()
	clusterName := "unknown"
	if conf, err := config.GetConfig(); err == nil {
		clusterName = conf.ClusterName
	}

	for resourceObj := range inChan {
		rc.storeDevice(lctx, resourceObj, clusterName, store)
	}
	log.Debugf("New cache map : %v", store)
}

func (rc *ResourceCache) storeDevice(lctx *lmctx.LMContext, resourceObj *models.Device, clusterName string, store *Store) bool {
	log := lmlog.Logger(lctx)
	if resourceObj == nil ||
		resourceObj.DeviceType != 8 ||
		util.GetResourcePropertyValue(resourceObj, "auto.clustername") != clusterName ||
		util.GetResourcePropertyValue(resourceObj, constants.K8sResourceDeletedOnPropertyKey) != "" {
		return false
	}

	log.Debugf("Accumulating Resource %v: %v", resourceObj.ID, resourceObj.DisplayName)
	rt, er := util.GetResourceType(resourceObj)
	if er != nil {
		log.Errorf("ResourceType cannot be determinied using resource object : %s", er)

		return false
	}
	key := types.ResourceName{
		Name:     util.GetResourcePropertyValue(resourceObj, constants.K8sResourceNamePropertyKey),
		Resource: rt,
	}

	meta, err := util.GetResourceMetaFromResource(resourceObj)
	if err != nil {
		log.Debugf("Cannot get resource meta to store in cache: %v", resourceObj)

		return false
	}

	// ignore deleted category resources
	if meta.HasSysCategory(rt.GetDeletedCategory()) {
		return false
	}

	store.Set(key, meta)
	return true
}

func (rc *ResourceCache) accumulateDeviceGroupCache(lctx *lmctx.LMContext, inChan <-chan DeviceGroupData, store *Store, finished chan<- bool) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		finished <- true
		close(finished)
		log.Infof("ResourceGroup accumulation completed in %v seconds", time.Since(start).Seconds())
	}()
	for resourceGroup := range inChan {
		if resourceGroup.ResourceName == "" {
			continue
		}
		log.Debugf("Accumulating ResourceGroup %v: %v", resourceGroup.ResourceName, resourceGroup.NamespaceName)
		key := types.ResourceName{
			Name:     resourceGroup.ResourceName,
			Resource: enums.Namespaces,
		}
		customProps := make(map[string]string)
		for _, nv := range resourceGroup.CustomProps {
			customProps[*nv.Name] = *nv.Value
		}

		meta := types.ResourceMeta{ // nolint: exhaustivestruct
			Container: fmt.Sprintf("%d", resourceGroup.NamespaceName),
			LMID:      resourceGroup.LMID,
			Labels:    customProps,
		}
		store.Set(key, meta)
	}
	log.Debugf("New cache map ResourceGroup: %v", store)
}

func (rc *ResourceCache) resetCacheStore(m *Store) {
	rc.rwm.Lock()
	defer rc.rwm.Unlock()
	rc.store = m
	go func() {
		rc.hookrwm.RLock()
		defer rc.hookrwm.RUnlock()
		for _, hook := range rc.hooks {
			// run hook on existing items
			for _, item := range rc.List() {
				if hook.Predicate(types.CacheSet, item.K, item.V) {
					hook.Hook(item.K, item.V)
				}
			}
		}
	}()
}

func (rc *ResourceCache) getAllGroups(lctx *lmctx.LMContext, grpid int32, outChan chan<- int32, groupChan chan<- DeviceGroupData, depth int) {
	log := lmlog.Logger(lctx)

	if depth < 1 {
		return
	}
	g, err := rc.getDeviceGroupByID(lctx, grpid)
	if err != nil {
		log.Errorf("Failed to fetch group with id: %v", grpid)

		return
	}
	// errRegex := regexp.MustCompile(`ns(?P<code>\d+)`)
	for _, sg := range g.Payload.SubGroups {
		// Custom props for subgroups are not stored in cache yet, as no use case, and fetching all groups will be time consuming and increases lm calls
		groupChan <- DeviceGroupData{
			ResourceName:  sg.Name,
			NamespaceName: g.Payload.ID,
			LMID:          sg.ID,
		}

		if sg.Name == "_deleted" {
			continue
		}

		log.Tracef("Taking group: %v of parent %d", sg.Name, grpid)
		outChan <- sg.ID
		rc.getAllGroups(lctx, sg.ID, outChan, groupChan, depth-1)
	}
}

func (rc *ResourceCache) getDeviceGroupByID(lctx *lmctx.LMContext, grpid int32) (*lm.GetDeviceGroupByIDOK, error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(grpid)
	command := rc.GetResourceGroupByIDCommand(clctx, params)
	resp, err := rc.SendReceive(clctx, command)
	if err != nil {
		return nil, err
	}

	return resp.(*lm.GetDeviceGroupByIDOK), err
}

// Set adds entry into cache map
func (rc *ResourceCache) Set(lctx *lmctx.LMContext, name types.ResourceName, meta types.ResourceMeta) bool {
	log := lmlog.Logger(lctx)
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log.Tracef("Deleting previous entry if any...")
	rc.unsetInternal(lctx, name, meta.Container)
	log.Tracef("Setting cache entry %v: %v", name, meta)
	ok := rc.store.Set(name, meta)
	if ok {
		go func() {
			rc.hookrwm.RLock()
			defer rc.hookrwm.RUnlock()
			for _, hook := range rc.hooks {
				if hook.Predicate(types.CacheSet, name, meta) {
					hook.Hook(name, meta)
				}
			}
		}()
	}
	return ok
}

// Exists checks entry into cache map
func (rc *ResourceCache) Exists(lctx *lmctx.LMContext, name types.ResourceName, container string, softRefresh bool) (types.ResourceMeta, bool) {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %v: %v", name, container)
	meta, ok := rc.store.Exists(lctx, name, container)
	if !ok && softRefresh {
		rc.SoftRefresh(lctx, container)
	}
	if softRefresh {
		return rc.store.Exists(lctx, name, container)
	}
	return meta, ok
}

// Get checks entry into cache map
func (rc *ResourceCache) Get(lctx *lmctx.LMContext, name types.ResourceName) ([]types.ResourceMeta, bool) {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)

	return rc.store.Get(lctx, name)
}

// Unset checks entry into cache map
func (rc *ResourceCache) Unset(lctx *lmctx.LMContext, name types.ResourceName, container string) bool {
	rc.rwm.RLock()
	defer rc.rwm.RUnlock()
	meta, ok := rc.unsetInternal(lctx, name, container)
	if ok {
		go func() {
			rc.hookrwm.RLock()
			defer rc.hookrwm.RUnlock()
			for _, hook := range rc.hooks {
				if hook.Predicate(types.CacheUnset, name, meta) {
					hook.Hook(name, meta)
				}
			}
		}()
	}
	return ok
}

func (rc *ResourceCache) unsetInternal(lctx *lmctx.LMContext, name types.ResourceName, container string) (types.ResourceMeta, bool) {
	log := lmlog.Logger(lctx)
	/*// Special handling for resource groups, do not add another.
	// when parent resource group is deleted its all child hierarchy gets deleted from portal, hence
	if name.Resource == enums.Namespaces {
		if exists, ok := rc.store.Exists(lctx, name, container); ok {
			list := rc.store.ListWithFilter(func(k cache.ResourceName, v cache.ResourceMeta) bool {
				return v.Container == exists.Container ||
					(k.Resource.IsNamespaceScopedResource() && v.Container == name.GroupName)
			})
			log.Infof("Removing dependent (%d) container cache entries from cache of %s", len(list), container)
			for _, item := range list {
				rc.store.Unset(item.K, item.V.Container)
			}
		}
	}*/
	log.Tracef("Deleting cache entry %v: %v", name, container)

	return rc.store.Unset(name, container)
}

// Load loads cache from configmaps
func (rc *ResourceCache) Load(lctx *lmctx.LMContext) error {
	cmList, err := rc.listAllCacheConfigMaps()
	if err != nil {
		return err
	}
	if cmList.Size() == 0 {
		return fmt.Errorf("no config maps found")
	}
	selectedDumpID := rc.selectDumpID(lctx, cmList)
	tmpCache := NewStore()
	err2 := rc.populateCacheStore(lctx, cmList, selectedDumpID, tmpCache)
	if err2 != nil {
		return err2
	}
	if tmpCache.Size() == 0 {
		return fmt.Errorf("no cache data found in present configmaps")
	}
	rc.resetCacheStore(tmpCache)
	rc.dumpID = selectedDumpID
	rc.updateIncrementalCache(lctx)

	return nil
}

func (rc *ResourceCache) populateCacheStore(lctx *lmctx.LMContext, cmList *corev1.ConfigMapList, selectedDumpID int64, tmpCache *Store) error {
	log := lmlog.Logger(lctx)
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 != nil || (selectedDumpID != -1 && dumpID != selectedDumpID) {
			continue
		}
		m := make(map[types.ResourceName][]types.ResourceMeta)
		er := json.Unmarshal([]byte(cm.Data["cache"]), &m)
		if er != nil {
			log.Errorf("Failed to parse stored configmap [%s]: %s", cm.Name, er)

			return er
		}
		tmpCache.AddAll(m)
	}

	return nil
}

func (rc *ResourceCache) selectDumpID(lctx *lmctx.LMContext, cmList *corev1.ConfigMapList) int64 {
	log := lmlog.Logger(lctx)
	// Select latest dump if multiple dump exists
	selectedDumpID := int64(-1)
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 == nil && dumpID > selectedDumpID {
			selectedDumpID = dumpID
		}
	}
	if selectedDumpID == -1 {
		log.Warn("No dumpID found in any of the listed configMap, loading all configmaps into cache (may go wrong if multiple dumps are present)")
	}

	return selectedDumpID
}

func (rc *ResourceCache) listAllCacheConfigMaps() (*corev1.ConfigMapList, error) {
	ns, err := config.GetWatchConfig("namespace")
	if err != nil {
		ns = "argus"
	}
	cmList, err := config.GetClientSet().CoreV1().ConfigMaps(ns).List(metav1.ListOptions{
		LabelSelector: labels.Set{"argus": "cache"}.String(),
	})
	if err != nil {
		return cmList, fmt.Errorf("failed to list all configmaps with selector \"argus=cache\": %w", err)
	}
	return cmList, nil
}

// Save saves cache to configmaps
func (rc *ResourceCache) Save() error {
	rc.flushMU.Lock()
	defer rc.flushMU.Unlock()
	if rc.store.Size() == 0 {
		logrus.Tracef("store is empty so not storing it")

		return fmt.Errorf("store is empty hence not storing to cm")
	}

	chunks, err := getChunks(rc.store.getMap())
	if err != nil {
		logrus.Errorf("Failed to marshal cache map to json string or failed to split into chunks %v", err)

		return err
	}
	ns, err := config.GetWatchConfig("namespace")
	if err != nil {
		ns = "argus"
	}
	dumpID := time.Now().Unix()
	dumpIDStr := fmt.Sprintf("%v", dumpID)
	for idx, chunk := range chunks {
		m := map[string]string{"cache": fmt.Sprintf("%s", chunk)}
		cm := &corev1.ConfigMap{ // nolint: exhaustivestruct
			TypeMeta: metav1.TypeMeta{}, // nolint: exhaustivestruct
			ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
				Name:      fmt.Sprintf("cache-%v-%v", dumpID, idx),
				Namespace: ns,
				Labels: map[string]string{
					"argus":       "cache",
					"dumpID":      dumpIDStr,
					"chunkNumber": fmt.Sprintf("%v", idx),
					"version":     Version,
				},
				Annotations: map[string]string{
					"content_size": fmt.Sprintf("%v", rc.store.Size()),
				},
			},
			Data:       m,
			BinaryData: nil,
		}
		_, err1 := config.GetClientSet().CoreV1().ConfigMaps(ns).Create(cm)
		if err1 != nil {
			_, err2 := config.GetClientSet().CoreV1().ConfigMaps(ns).Update(cm)
			if err2 != nil {
				err3 := fmt.Errorf("failed to store cache chunk %v to cm: %w %v", idx, err1, err2)
				logrus.Errorf("%s", err3)

				return err3
			}
		}
	}
	// Delete previous cache configmaps
	err2 := config.GetClientSet().CoreV1().ConfigMaps(ns).DeleteCollection(
		&metav1.DeleteOptions{}, // nolint: exhaustivestruct
		metav1.ListOptions{ // nolint: exhaustivestruct
			LabelSelector: fmt.Sprintf("argus==cache,dumpID!=%s", dumpIDStr),
		})
	if err2 != nil {
		return fmt.Errorf("failed to delete previous cache configmaps: %w", err2)
	}

	return nil
}

func getChunks(m map[types.ResourceName][]types.ResourceMeta) ([][]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var result [][]byte
	if len(bytes) > CMMaxBytes {
		n := 1
		m1 := make(map[types.ResourceName][]types.ResourceMeta)
		m2 := make(map[types.ResourceName][]types.ResourceMeta)
		for k, v := range m {
			if n%2 == 0 {
				m1[k] = v
			} else {
				m2[k] = v
			}
			n++
		}
		ch1, err2 := getChunks(m1)
		if err2 != nil {
			return result, err2
		}

		result = append(result, ch1...)
		ch2, err3 := getChunks(m2)
		if err3 != nil {
			return result, err3
		}

		result = append(result, ch2...)

		return result, nil
	}
	result = append(result, bytes)

	return result, nil
}

// List returns slices of all data present in cache - For ex: cache for periodic dangling resource etc uses this
func (rc *ResourceCache) List() []types.IterItem {
	return rc.store.List()
}

// ListWithFilter returns slices of all data present in cache - For ex: cache for periodic dangling resource etc uses this
func (rc *ResourceCache) ListWithFilter(f func(k types.ResourceName, v types.ResourceMeta) bool) []types.IterItem {
	return rc.store.ListWithFilter(f)
}

// UnsetLMID unsets value using id only
// WARN: This is not an o(1) map operation hence do not use until necessary, this is heavy operation to find cache entry with id
func (rc *ResourceCache) UnsetLMID(lctx *lmctx.LMContext, rt enums.ResourceType, id int32) bool {
	k, m := rc.store.getInternalUsingLMID(rt, id)
	return rc.Unset(lctx, k, m.Container)
}

func (rc *ResourceCache) updateIncrementalCache(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Update incremental cache that might not have dumped into configmap cache..")
	conf, err := config.GetConfig()
	if err != nil {
		log.Warnf("Failed to get config")
		return
	}
	for {
		clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
		params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
		params.SetID(*conf.ResourceContainerGroupID)
		// keeping 10 seconds less to cover corner cases - wherein cache was just dumped and objects updated
		// https://www.logicmonitor.com/support/rest-api-developers-guide/v1/resources/get-resources
		filter := fmt.Sprintf("updatedOn>:\"%d\"", rc.dumpID-10) // nolint: gomnd
		params.SetFilter(&filter)
		command := rc.GetImmediateResourceListByResourceGroupIDCommand(clctx, params)
		resp, err := rc.SendReceive(clctx, command)
		if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusTooManyRequests {
			log.Warnf("Rate limits reached, retrying in %f seconds...", rateLimitBackoffTime.Seconds())
			time.Sleep(rateLimitBackoffTime)

			continue
		}
		count := 0
		if resp != nil && util.GetHTTPStatusCodeFromLMSDKError(resp.(*lm.GetImmediateDeviceListByDeviceGroupIDOK)) == http.StatusOK {
			for _, resource := range resp.(*lm.GetImmediateDeviceListByDeviceGroupIDOK).Payload.Items {
				if rc.storeDevice(lctx, resource, conf.ClusterName, rc.store) {
					count++
				}
			}
		}
		log.Infof("Updated incremental cache update with %d entries", count)
		break
	}
}

func (rc *ResourceCache) SoftRefresh(lctx *lmctx.LMContext, container string) {
	conf, err := config.GetConfig()
	if err != nil {
		return
	}
	if t, ok := rc.getSoftRefreshLastTime(container); ok && t.Before(time.Now().Add(-time.Minute)) {
		list, ok := rc.Get(lctx, types.ResourceName{Name: container, Resource: enums.Namespaces})
		if ok {
			for _, meta := range list {
				resp, err := rc.getDevices(nil, meta.LMID)
				if err != nil && resp != nil {
					for _, resource := range resp.Payload.Items {
						rc.storeDevice(lctx, resource, conf.ClusterName, rc.store)
					}
				}
			}
		}
	}
}

func (rc *ResourceCache) getSoftRefreshLastTime(key string) (time.Time, bool) {
	rc.softRefreshMu.Lock()
	defer rc.softRefreshMu.Unlock()
	if val, ok := rc.softRefreshLast[key]; ok {
		return val, ok
	}

	return time.Time{}, false
}
