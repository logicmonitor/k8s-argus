package resourcecache

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

func (rc *ResourceCache) rebuildCache() {
	rc.rebuildMutex.Lock()
	defer rc.rebuildMutex.Unlock()
	defer metrics.ObserveTime(metrics.StartTimeObserver(metrics.CacheBuilderSummary))
	debugID := util.GetShortUUID()
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"debug_id": debugID}))
	log := lmlog.Logger(lctx)
	log.Infof("Resource cache fetching resources")
	resources, err := rc.getAllResources(lctx)
	if resources == nil || err != nil {
		log.Errorf("Failed to fetch resources: %s", err)
	} else {
		log.Debugf("Resync cache map")
		rc.resetCacheStore(resources)
		log.Debugf("Resync cache done")
	}

	if err := rc.Save(lctx); err != nil {
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
	clusterGroupID, err := util.GetClusterGroupID(lctx, rc.LMRequester)
	if err != nil {
		log.Error(err.Error())

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

	go rc.fetchGroupDevices(lctx, grpIDChan, resourceChan)

	grpIDChan <- clusterGroupID
	if conf, err := config.GetConfig(lctx); err == nil {
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
	conf, err := config.GetConfig(lctx)
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
	conf, err := config.GetConfig(lctx)
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

func (rc *ResourceCache) fetchGroupDevices(lctx *lmctx.LMContext, inChan <-chan int32, outChan chan<- *models.Device) {
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
	if conf, err := config.GetConfig(lctx); err == nil {
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

	store.Set(lctx, key, meta)
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
		store.Set(lctx, key, meta)
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
	conf, err := config.GetConfig(lctx)
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
