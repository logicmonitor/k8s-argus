package devicecache

import (
	"strings"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

// DeviceCache to maintain a device cache to calculate delta between device presence on server and on cluster
type DeviceCache struct {
	store       map[string]interface{}
	rwm         sync.RWMutex
	rsyncPeriod time.Duration
	base        *types.Base
}

// NewDeviceCache create new DeviceCache object
func NewDeviceCache(b *types.Base, rp time.Duration) *DeviceCache {
	dc := &DeviceCache{
		store:       make(map[string]interface{}),
		base:        b,
		rsyncPeriod: rp,
		rwm:         sync.RWMutex{},
	}

	return dc
}

// Run start a goroutine to rsync cache periodically with server
func (dc *DeviceCache) Run() {
	go func(dcache *DeviceCache) {
		for {
			start := time.Now()
			logrus.Debugf("Device cache fetching devices")
			devices := dcache.getAllDevices(dcache.base)
			if devices == nil {
				logrus.Errorf("Failed to fetch devices")
			} else {
				logrus.Debugf("Resync cache map")
				dcache.rsyncCache(devices)
				logrus.Debugf("Resync cache done")
			}
			end := time.Now()
			logrus.Infof("Old Cache synced in %v", end.Sub(start).Milliseconds())
			time.Sleep(dcache.rsyncPeriod)
		}
	}(dc)
}

func (dc *DeviceCache) getAllDevices(b *types.Base) map[string]interface{} {
	cgrpid := b.Config.ClusterGroupID

	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(cgrpid)

	g, err := b.LMClient.LM.GetDeviceGroupByID(params)
	if err != nil {
		logrus.Errorf("Error while fetching cluster device group %v", err)

		return nil
	}

	clusterGroupName := util.ClusterGroupName(b.Config.ClusterName)
	clusterGroupID := int32(0)

	for _, sg := range g.Payload.SubGroups {
		if sg.Name == clusterGroupName {
			clusterGroupID = sg.ID

			break
		}
	}

	if clusterGroupID == 0 {
		logrus.Errorf("No Cluster group found")

		return nil
	}

	grps := dc.getAllGroups(b, clusterGroupID)
	grps = append(grps, clusterGroupID)

	logrus.Debugf("all groups: %#v", grps)
	m := make(map[string]interface{})
	for _, gid := range grps {
		params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
		params.SetID(gid)
		resp, err := b.LMClient.LM.GetImmediateDeviceListByDeviceGroupID(params)
		if err != nil {
			continue
		}
		for _, device := range resp.Payload.Items {
			m[dc.getFullDisplayName(device)] = true
		}
	}

	return m
}

func getResourceTypeFromSystemCateogries(category string) enums.ResourceType {
	if strings.Contains(category, constants.PodCategory) {
		return enums.Pods
	}
	if strings.Contains(category, constants.DeploymentCategory) {
		return enums.Deployments
	}
	if strings.Contains(category, constants.ServiceCategory) {
		return enums.Services
	}
	if strings.Contains(category, constants.NodeCategory) {
		return enums.Nodes
	}
	if strings.Contains(category, constants.HorizontalPodAutoscalerCategory) {
		return enums.Hpas
	}

	return enums.Unknown
}

func (dc *DeviceCache) getFullDisplayName(device *models.Device) string {
	systemCategories := util.GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)

	return util.GetFullDisplayName(device, getResourceTypeFromSystemCateogries(systemCategories), dc.base.Config.ClusterName)
}

func (dc *DeviceCache) rsyncCache(m map[string]interface{}) {
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	dc.store = m
}

func (dc *DeviceCache) getAllGroups(b *types.Base, grpID int32) []int32 {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(grpID)
	g, err := b.LMClient.LM.GetDeviceGroupByID(params)
	if err != nil {
		logrus.Errorf("Failed to fetch group with id: %v", grpID)

		return []int32{}
	}
	var subGroups []int32

	for _, sg := range g.Payload.SubGroups {
		if sg.Name == "_deleted" {
			continue
		}
		logrus.Debugf("Taking group: %v", sg.Name)
		gps := dc.getAllGroups(b, sg.ID)
		subGroups = append(subGroups, gps...)
	}

	subGroups = append(subGroups, grpID)

	return subGroups
}

// Set adds entry into cache map
func (dc *DeviceCache) Set(name string) bool {
	logrus.Debugf("Setting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	dc.store[name] = true

	return true
}

// Exists checks entry into cache map
func (dc *DeviceCache) Exists(name string) bool {
	logrus.Debugf("Checking cache entry %s", name)
	dc.rwm.RLock()
	defer dc.rwm.RUnlock()
	_, ok := dc.store[name]

	return ok
}

// Unset checks entry into cache map
func (dc *DeviceCache) Unset(name string) bool {
	logrus.Debugf("Deleting cache entry %s", name)
	dc.rwm.Lock()
	defer dc.rwm.Unlock()
	delete(dc.store, name)

	return true
}
