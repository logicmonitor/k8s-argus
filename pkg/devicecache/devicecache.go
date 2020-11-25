package devicecache

import (
	"strings"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"github.com/vkumbhar94/lm-sdk-go/client/lm"
	"github.com/vkumbhar94/lm-sdk-go/models"
)

// DeviceCache to maintain a device cache to calcuate delta between device presence on server and on cluster
type DeviceCache struct {
	store        map[string]interface{}
	rwm          sync.RWMutex
	resyncPeriod int
	base         *types.Base
}

// NewDeviceCache create new DeviceCache object
func NewDeviceCache(b *types.Base, rp int) *DeviceCache {
	dc := &DeviceCache{
		store:        make(map[string]interface{}),
		base:         b,
		resyncPeriod: rp,
	}
	return dc
}

// Run start a goroutine to resync cache periodically with server
func (dc *DeviceCache) Run() {
	go func(dcache *DeviceCache) {
		for {
			time.Sleep(time.Duration(dcache.resyncPeriod) * time.Minute)
			log.Debugf("Device cache fetching devices")
			devices := dcache.getAllDevices(dcache.base)
			if devices == nil {
				log.Errorf("Failed to fetch devices")
			} else {
				log.Debugf("Resync cache map")
				dcache.resyncCache(devices)
				log.Debugf("Resync cache done")
			}
		}
	}(dc)
}

func (dc *DeviceCache) getAllDevices(b *types.Base) map[string]interface{} {
	cgrpid := b.Config.ClusterGroupID

	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(cgrpid)

	g, err := b.LMClient.LM.GetDeviceGroupByID(params)
	if err != nil {
		log.Errorf("Error while fetching cluster device group %v", err)
		return nil
	}

	clusterGroupName := constants.ClusterDeviceGroupPrefix + b.Config.ClusterName
	clusterGroupID := int32(0)

	for _, sg := range g.Payload.SubGroups {
		if sg.Name == clusterGroupName {
			clusterGroupID = sg.ID
			break
		}
	}

	if clusterGroupID == 0 {
		log.Errorf("No Cluster group found")
		return nil
	}

	grps := dc.getAllGroups(b, clusterGroupID)
	grps = append(grps, clusterGroupID)

	log.Debugf("all groups: %#v", grps)
	m := make(map[string]interface{})
	for _, gid := range grps {
		params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
		params.SetID(gid)
		resp, err := b.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
		if err != nil {
			continue
		}
		for _, device := range resp.Payload.Items {
			fullname := dc.getFullDisplayName(device)
			log.Debugf("devicecase added entry for - %v", fullname)
			m[fullname] = true
		}
	}
	return m
}

func (dc *DeviceCache) getFullDisplayName(device *models.Device) string {
	syscategory := util.GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)
	return util.GetFullDisplayName(device, getResourceTypeFromSystemCateogries(syscategory), dc.base.Config.ClusterName)
}

func getResourceTypeFromSystemCateogries(category string) string {
	if strings.Contains(category, constants.PodCategory) {
		return constants.Pods
	}
	if strings.Contains(category, constants.DeploymentCategory) {
		return constants.Deployments
	}
	if strings.Contains(category, constants.ServiceCategory) {
		return constants.Services
	}
	if strings.Contains(category, constants.NodeCategory) {
		return constants.Nodes
	}
	if strings.Contains(category, constants.HorizontalPodAutoscalerCategory) {
		return constants.HorizontalPodAutoScalers
	}
	return ""
}

func (dc *DeviceCache) resyncCache(m map[string]interface{}) {
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	dc.store = m
}

func (dc *DeviceCache) getAllGroups(b *types.Base, grpid int32) []int32 {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(grpid)
	g, err := b.LMClient.LM.GetDeviceGroupByID(params)
	if err != nil {
		log.Errorf("Failed to fetch group with id: %v", grpid)
		return []int32{}
	}
	subGroups := []int32{}
	for _, sg := range g.Payload.SubGroups {
		if sg.Name == "_deleted" {
			continue
		}
		log.Debugf("Taking group: %v", sg.Name)
		gps := dc.getAllGroups(b, sg.ID)
		subGroups = append(subGroups, gps...)
	}
	subGroups = append(subGroups, grpid)
	return subGroups
}

// Set adds entry into cache map
func (dc *DeviceCache) Set(name string) bool {
	log.Debugf("Setting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	dc.store[name] = true
	return true
}

// Exists checks entry into cache map
func (dc *DeviceCache) Exists(name string) bool {
	log.Debugf("Checking cache entry %s", name)
	dc.rwm.RLock()
	defer dc.rwm.RUnlock()
	_, ok := dc.store[name]
	return ok
}

// Unset checks entry into cache map
func (dc *DeviceCache) Unset(name string) bool {
	log.Debugf("Deleting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	delete(dc.store, name)
	return true
}
