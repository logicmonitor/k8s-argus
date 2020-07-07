package devicecache

import (
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	log "github.com/sirupsen/logrus"
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
func (dc *DeviceCache) Run() error {
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
	return nil
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
			m[*device.DisplayName] = true
		}
	}
	return m
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
