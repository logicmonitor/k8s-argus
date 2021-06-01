package devicecache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
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

// ResourceCache to maintain a device cache to calculate delta between device presence on server and on cluster
type ResourceCache struct {
	store         *Store
	rwm           sync.RWMutex
	rebuildMutex  sync.Mutex
	resyncPeriod  time.Duration
	flushTimeToCM time.Duration
	flushMU       sync.Mutex
	base          *types.Base
	clusterGrpID  int32
	stateLoaded   bool
}

// NewResourceCache create new ResourceCache object
func NewResourceCache(b *types.Base, rp time.Duration) *ResourceCache {
	resourceCache := &ResourceCache{
		store:         NewStore(),
		rwm:           sync.RWMutex{},
		resyncPeriod:  rp,
		flushTimeToCM: 1 * time.Minute,
		flushMU:       sync.Mutex{},
		base:          b,
		clusterGrpID:  -1,
		stateLoaded:   false,
		rebuildMutex:  sync.Mutex{},
	}

	if err := resourceCache.Load(); err != nil {
		logrus.Warnf("ResourceCache is not loaded from its last state stored in configmaps, it may take more time to rebuild cache on first run. Error by Load: %s", err)
	} else {
		resourceCache.stateLoaded = true
		logrus.Info("ResourceCache loaded from configmaps")
	}

	return resourceCache
}

// Run start a goroutine to resync cache periodically with server
func (resourceCache *ResourceCache) Run() {
	initialised := make(chan bool)
	go resourceCache.AutoCacheBuilder(initialised)

	// Periodically flush cache to configmap for recovery after application restart or pod restart
	go resourceCache.CacheToConfigMapDumper()

	// Wait for first initialization
	<-initialised
}

func (resourceCache *ResourceCache) AutoCacheBuilder(initialised chan<- bool) {
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "resource",
			Subsystem:   "cache",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "auto"},
		})
	t := time.NewTicker(resourceCache.resyncPeriod)

	if !resourceCache.stateLoaded {
		resourceCache.rebuildCache(gauge)
		initialised <- true
	} else {
		initialised <- true
	}
	close(initialised)

	for {
		<-t.C
		resourceCache.rebuildCache(gauge)
	}
}

func (resourceCache *ResourceCache) CacheToConfigMapDumper() {
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
	}(resourceCache)
}

// Rebuild graceful rebuild cache
func (resourceCache *ResourceCache) Rebuild(lctx *lmctx.LMContext) {
	gauge := promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace:   "cache",
			Subsystem:   "resource",
			Name:        "rebuild_time",
			Help:        "Time taken to rebuild cache",
			ConstLabels: prometheus.Labels{"run": "graceful"},
		})
	resourceCache.rebuildCache(gauge)
}

func (resourceCache *ResourceCache) rebuildCache(gauge prometheus.Gauge) {
	resourceCache.rebuildMutex.Lock()
	defer resourceCache.rebuildMutex.Unlock()
	defer promperf.Duration(promperf.Track(gauge))
	debugID := util.GetShortUUID()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"debug_id": debugID}))
	log := lmlog.Logger(lctx)
	log.Infof("Device cache fetching devices")
	devices := resourceCache.getAllDevices(lctx, resourceCache.base)
	if devices == nil {
		log.Errorf("Failed to fetch devices")
	} else {
		log.Debugf("Resync cache map")
		resourceCache.resetCacheStore(devices)
		log.Debugf("Resync cache done")
	}

	if err := resourceCache.Save(); err != nil {
		log.Errorf("cache to cm failed: %s", err)
	} else {
		log.Infof("Cache size: %v", resourceCache.store.Size())
	}
}

// DeviceGroupData schema to send device groups from api response
type DeviceGroupData struct {
	ResourceName  string
	NamespaceName string
	LMID          int32
}

func (resourceCache *ResourceCache) getAllDevices(lctx *lmctx.LMContext, b *types.Base) *Store {
	log := lmlog.Logger(lctx)
	clusterGroupID := util.GetClusterGroupID(lctx, b.LMClient)

	if clusterGroupID == -1 {
		log.Errorf("No Cluster group found")

		return nil
	}
	m := NewStore()
	deviceChan := make(chan *models.Device)
	deviceGroupChan := make(chan DeviceGroupData)
	deviceFinished := make(chan bool)
	deviceGroupFinished := make(chan bool)

	go resourceCache.accumulateDeviceCache(lctx, deviceChan, m, deviceFinished)
	go resourceCache.accumulateDeviceGroupCache(lctx, deviceGroupChan, m, deviceGroupFinished)

	grpIDChan := make(chan int32)

	go resourceCache.fetchGroupDevices(lctx, b, grpIDChan, deviceChan)

	grpIDChan <- clusterGroupID
	start := time.Now()
	resourceCache.getAllGroups(lctx, b, clusterGroupID, grpIDChan, deviceGroupChan, 2)
	log.Infof("Device group fetched in %v seconds", time.Since(start).Seconds())
	close(grpIDChan)
	close(deviceGroupChan)
	<-deviceGroupFinished
	<-deviceFinished

	return m
}

func getDevices(b *types.Base, grpID int32) (*lm.GetImmediateDeviceListByDeviceGroupIDOK, error) {
	params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
	params.SetID(grpID)
	resp, err := b.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (resourceCache *ResourceCache) fetchGroupDevices(lctx *lmctx.LMContext, b *types.Base, inChan <-chan int32, outChan chan<- *models.Device) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		close(outChan)
		log.Infof("Device fetch completed in %v seconds", time.Since(start).Seconds())
	}()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go resourceCache.ResourceGroupProcessor(inChan, log, b, outChan, &wg)
	}
	wg.Wait()
}

func (resourceCache *ResourceCache) ResourceGroupProcessor(inChan <-chan int32, log *logrus.Entry, b *types.Base, outChan chan<- *models.Device, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for grpID := range inChan {
		log.Debugf("Fetching devices of group %v", grpID)
		for {
			resp, err := getDevices(b, grpID)
			if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusTooManyRequests {
				log.Warnf("Rate limits reached, retrying in 10 seconds...")
				time.Sleep(rateLimitBackoffTime)

				continue
			}
			if resp != nil && util.GetHTTPStatusCodeFromLMSDKError(resp) == http.StatusOK {
				for _, device := range resp.Payload.Items {
					outChan <- device
				}

				break
			} else {
				log.Warnf("fetch devices for %v failed with %v", grpID, resp)

				break
			}
		}
	}
}

func (resourceCache *ResourceCache) accumulateDeviceCache(lctx *lmctx.LMContext, inChan <-chan *models.Device, store *Store, finished chan<- bool) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		finished <- true
		close(finished)
		log.Infof("Device accumulation completed in %v seconds", time.Since(start).Seconds())
	}()
	for deviceObj := range inChan {
		if deviceObj == nil {
			continue
		}
		if util.GetPropertyValue(deviceObj, constants.K8sResourceDeletedOnPropertyKey) != "" {
			continue
		}
		log.Debugf("Accumulating Resource %v: %v", deviceObj.ID, deviceObj.DisplayName)
		rt, er := util.GetResourceType(deviceObj)
		if er != nil {
			log.Errorf("ResourceType cannot be determinied using device object : %s", er)

			continue
		}
		key := cache.ResourceName{
			Name:     util.GetPropertyValue(deviceObj, constants.K8sDeviceNamePropertyKey),
			Resource: rt,
		}

		meta, err := util.GetResourceMetaFromDevice(deviceObj)
		if err != nil {
			log.Debugf("Cannot get device meta to store in cache: %v", deviceObj)

			continue
		}

		// ignore deleted category resources
		if meta.HasSysCategory(rt.GetDeletedCategory()) {
			continue
		}

		store.Set(key, meta)
	}
	log.Infof("New cache map : %v", store)
}

func (resourceCache *ResourceCache) accumulateDeviceGroupCache(lctx *lmctx.LMContext, inChan <-chan DeviceGroupData, store *Store, finished chan<- bool) {
	log := lmlog.Logger(lctx)
	start := time.Now()
	defer func() {
		finished <- true
		close(finished)
		log.Infof("ResourceGroup accumulation completed in %v seconds", time.Since(start).Seconds())
	}()
	for deviceGroup := range inChan {
		if deviceGroup.ResourceName == "" || deviceGroup.NamespaceName == "" {
			continue
		}
		log.Debugf("Accumulating ResourceGroup %v: %v", deviceGroup.ResourceName, deviceGroup.NamespaceName)
		key := cache.ResourceName{
			Name:     fmt.Sprintf("ns-%s", deviceGroup.ResourceName),
			Resource: enums.Namespaces,
		}

		meta := cache.ResourceMeta{ // nolint: exhaustivestruct
			Container: deviceGroup.NamespaceName,
			LMID:      deviceGroup.LMID,
		}
		store.Set(key, meta)
	}
	log.Debugf("New cache map ResourceGroup: %v", store)
}

func (resourceCache *ResourceCache) resetCacheStore(m *Store) {
	resourceCache.rwm.Lock()
	defer resourceCache.rwm.Unlock()
	resourceCache.store = m
}

func (resourceCache *ResourceCache) getAllGroups(lctx *lmctx.LMContext, b *types.Base, grpid int32, outChan chan<- int32, groupChan chan<- DeviceGroupData, depth int) {
	log := lmlog.Logger(lctx)

	if depth < 1 {
		return
	}
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(grpid)
	g, err := b.LMClient.LM.GetDeviceGroupByID(params)
	if err != nil {
		log.Errorf("Failed to fetch group with id: %v", grpid)

		return
	}

	// errRegex := regexp.MustCompile(`ns(?P<code>\d+)`)
	for _, sg := range g.Payload.SubGroups {
		if sg.Name == "_deleted" {
			continue
		}

		/*matches := errRegex.FindStringSubmatch(sg.Name)
				if len(matches) > 0 {
					code, err := strconv.Atoi(matches[1])
					if err == nil && code > 50 {
						log.Warnf("Temporary debug logic, remove it %v : %v", sg.ID, sg.Name)

		continue
					}
				}*/

		log.Tracef("Taking group: %v", sg.Name)
		outChan <- sg.ID
		groupChan <- DeviceGroupData{
			ResourceName:  sg.Name,
			NamespaceName: *g.Payload.Name,
			LMID:          sg.ID,
		}
		resourceCache.getAllGroups(lctx, b, sg.ID, outChan, groupChan, depth-1)
	}
}

// Set adds entry into cache map
func (resourceCache *ResourceCache) Set(name cache.ResourceName, meta cache.ResourceMeta) bool {
	resourceCache.rwm.RLock()
	defer resourceCache.rwm.RUnlock()
	logrus.Tracef("Setting cache entry %v", name)
	resourceCache.store.Set(name, meta)

	return true
}

// Exists checks entry into cache map
func (resourceCache *ResourceCache) Exists(lctx *lmctx.LMContext, name cache.ResourceName, namespace string) (cache.ResourceMeta, bool) {
	resourceCache.rwm.RLock()
	defer resourceCache.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Checking cache entry %v", name)

	return resourceCache.store.Exists(lctx, name, namespace)
}

// Get checks entry into cache map
func (resourceCache *ResourceCache) Get(lctx *lmctx.LMContext, name cache.ResourceName) ([]cache.ResourceMeta, bool) {
	resourceCache.rwm.RLock()
	defer resourceCache.rwm.RUnlock()
	log := lmlog.Logger(lctx)
	log.Tracef("Get cache entry list %v", name)

	return resourceCache.store.Get(lctx, name)
}

// Unset checks entry into cache map
func (resourceCache *ResourceCache) Unset(name cache.ResourceName, namespace string) bool {
	resourceCache.rwm.RLock()
	defer resourceCache.rwm.RUnlock()
	logrus.Tracef("Deleting cache entry %v", name)

	return resourceCache.store.Unset(name, namespace)
}

// Load loads cache from configmaps
func (resourceCache *ResourceCache) Load() error {
	cmList, err := resourceCache.listAllCacheConfigMaps()
	if err != nil {
		return err
	}
	if cmList.Size() == 0 {
		return fmt.Errorf("no config maps found")
	}
	selectedDumpID := resourceCache.selectDumpID(cmList)
	tmpCache := NewStore()
	err2 := resourceCache.populateCacheStore(cmList, selectedDumpID, tmpCache)
	if err2 != nil {
		return err2
	}
	if tmpCache.Size() == 0 {
		return fmt.Errorf("no cache data found in present configmaps")
	}
	resourceCache.resetCacheStore(tmpCache)

	return nil
}

func (resourceCache *ResourceCache) populateCacheStore(cmList *corev1.ConfigMapList, selectedDumpID int64, tmpCache *Store) error {
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 != nil || (selectedDumpID != -1 && dumpID != selectedDumpID) {
			continue
		}
		m := make(map[cache.ResourceName][]cache.ResourceMeta)
		er := json.Unmarshal([]byte(cm.Data["cache"]), &m)
		if er != nil {
			logrus.Errorf("Failed to parse stored configmap [%s]: %s", cm.Name, er)

			return er
		}
		tmpCache.AddAll(m)
	}

	return nil
}

func (resourceCache *ResourceCache) selectDumpID(cmList *corev1.ConfigMapList) int64 {
	// Select latest dump if multiple dump exists
	selectedDumpID := int64(-1)
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 == nil && dumpID > selectedDumpID {
			selectedDumpID = dumpID
		}
	}
	if selectedDumpID == -1 {
		logrus.Warn("No dumpID found in any of the listed configMap, loading all configmaps into cache (may go wrong if multiple dumps are present)")
	}

	return selectedDumpID
}

func (resourceCache *ResourceCache) listAllCacheConfigMaps() (*corev1.ConfigMapList, error) {
	ns, err := config.GetWatchConfig("namespace")
	if err != nil {
		ns = "argus"
	}
	cmList, err := resourceCache.base.K8sClient.CoreV1().ConfigMaps(ns).List(metav1.ListOptions{
		LabelSelector: labels.Set{"argus": "cache"}.String(),
	})
	if err != nil {
		return cmList, fmt.Errorf("failed to list all configmaps with selector \"argus=cache\": %w", err)
	}
	return cmList, nil
}

// Save saves cache to configmaps
func (resourceCache *ResourceCache) Save() error {
	resourceCache.flushMU.Lock()
	defer resourceCache.flushMU.Unlock()
	if resourceCache.store.Size() == 0 {
		logrus.Tracef("store is empty so not storing it")

		return fmt.Errorf("store is empty hence not storing to cm")
	}
	chunks, err := getChunks(resourceCache.store.InternalMap)
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
					"content_size": fmt.Sprintf("%v", resourceCache.store.Size()),
				},
			},
			Data:       m,
			BinaryData: nil,
		}
		_, err1 := resourceCache.base.K8sClient.CoreV1().ConfigMaps(ns).Create(cm)
		if err1 != nil {
			_, err2 := resourceCache.base.K8sClient.CoreV1().ConfigMaps(ns).Update(cm)
			if err2 != nil {
				err3 := fmt.Errorf("failed to store cache chunk %v to cm: %w %v", idx, err1, err2)
				logrus.Errorf("%s", err3)

				return err3
			}
		}
	}
	// Delete previous cache configmaps
	err2 := resourceCache.base.K8sClient.CoreV1().ConfigMaps(ns).DeleteCollection(
		&metav1.DeleteOptions{}, // nolint: exhaustivestruct
		metav1.ListOptions{ // nolint: exhaustivestruct
			LabelSelector: fmt.Sprintf("argus==cache,dumpID!=%s", dumpIDStr),
		})
	if err2 != nil {
		return fmt.Errorf("failed to delete previous cache configmaps: %w", err2)
	}

	return nil
}

func getChunks(m map[cache.ResourceName][]cache.ResourceMeta) ([][]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var result [][]byte
	if len(bytes) > CMMaxBytes {
		n := 1
		m1 := make(map[cache.ResourceName][]cache.ResourceMeta)
		m2 := make(map[cache.ResourceName][]cache.ResourceMeta)
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

// List returns slices of all data present in cache - For ex: cache for periodic dangling device etc uses this
func (resourceCache *ResourceCache) List() []cache.IterItem {
	return resourceCache.store.List()
}

// UnsetLMID unsets value using id only
// WARN: This is not an o(1) map operation hence do not use until necessary, this is heavy operation to find cache entry with id
func (resourceCache *ResourceCache) UnsetLMID(rt enums.ResourceType, id int32) bool {
	return resourceCache.store.UnsetLMID(rt, id)
}
