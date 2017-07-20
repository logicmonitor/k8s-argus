package node

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/k8s-argus/argus/config"
	"github.com/logicmonitor/k8s-argus/argus/constants"
	lmv1 "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/runtime"
)

const (
	resource = "nodes"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct {
	LMClient  *lmv1.DefaultApi
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}

// Resource is a function that implements the Watcher interface.
func (w Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w Watcher) ObjType() runtime.Object {
	return &v1.Node{}
}

// AddFunc is a function that implements the Watcher interface.
func (w Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		if GetInternalAddress(node.Status.Addresses) != nil {
			log.Infof("Adding node %s", node.Name)
			w.addDevice(node)
		}
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		oldNode := oldObj.(*v1.Node)
		newNode := newObj.(*v1.Node)
		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		newInternalAddress := GetInternalAddress(oldNode.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			log.Infof("Adding node %s", newNode.Name)
			w.addDevice(newNode)
		} else if oldInternalAddress.Address != newInternalAddress.Address {
			// Covers the case when a node has been terminated (new ip doesn't exist)
			// and if a node needs to be added.
			filter := fmt.Sprintf("displayName:%s", oldNode.Name)
			restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
			if err != nil {
				log.Errorf("Failed searching for node %s: %s", oldNode.Name, restResponse.Errmsg)
			}
			if restResponse.Data.Total == 1 {
				id := restResponse.Data.Items[0].Id
				log.Infof("Updating node %s with id %d", oldNode.Name, id)
				w.updateDevice(newNode, id)
			}
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		filter := fmt.Sprintf("displayName:%s", node.Name)
		restResponse, _, err := w.LMClient.GetDeviceList("", -1, 0, filter)
		if err != nil {
			log.Errorf("Failed searching for node %s: %s", node.Name, restResponse.Errmsg)
		}
		if restResponse.Data.Total == 1 {
			id := restResponse.Data.Items[0].Id
			if w.Config.DeleteDevices {
				log.Infof("Deleting node %s with id %d", node.Name, id)
				restNullObjectResponse, _, err := w.LMClient.DeleteDevice(id)
				if err != nil {
					log.Printf("Failed to delete device with id %q: %s", id, restNullObjectResponse.Errmsg)
				}
			} else {
				log.Infof("Moving node %s with id %d to deleted group", node.Name, id)
				categories := constants.NodeDeletedCategory
				for k, v := range node.Labels {
					categories += "," + k + "=" + v

				}
				device := lmv1.RestDevice{
					CustomProperties: []lmv1.NameAndValue{
						{
							Name:  "system.categories",
							Value: categories,
						},
					},
				}
				restResponse, _, err := w.LMClient.PatchDeviceById(device, id, "replace", "customProperties")
				if err != nil {
					log.Errorf("Failed to patch node %s: %s", node.Name, restResponse.Errmsg)
				}
			}
		}
	}
}

func (w Watcher) addDevice(node *v1.Node) {
	device := w.makeDeviceObject(node)
	restResponse, _, err := w.LMClient.AddDevice(device, false)
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) updateDevice(node *v1.Node, id int32) {
	device := w.makeDeviceObject(node)
	restResponse, _, err := w.LMClient.UpdateDevice(device, id, "")
	if err != nil {
		log.Error(restResponse.Errmsg)
	}
}

func (w Watcher) makeDeviceObject(node *v1.Node) (device lmv1.RestDevice) {
	categories := constants.NodeCategory
	for k, v := range node.Labels {
		categories += "," + k + "=" + v

	}

	device = lmv1.RestDevice{
		Name:                 GetInternalAddress(node.Status.Addresses).Address,
		DisplayName:          node.Name,
		DisableAlerting:      w.Config.DisableAlerting,
		HostGroupIds:         "1",
		PreferredCollectorId: w.Config.PreferredCollector,
		CustomProperties: []lmv1.NameAndValue{
			{
				Name:  "system.categories",
				Value: categories,
			},
			{
				Name:  "auto.clustername",
				Value: w.Config.ClusterName,
			},
			{
				Name:  "auto.selflink",
				Value: node.SelfLink,
			},
			{
				Name:  "auto.name",
				Value: node.Name,
			},
			{
				Name:  "auto.uid",
				Value: string(node.UID),
			},
		},
	}

	return
}

// GetInternalAddress finds the node's internal address.
func GetInternalAddress(addresses []v1.NodeAddress) *v1.NodeAddress {
	for _, address := range addresses {
		if address.Type == v1.NodeInternalIP {
			return &address
		}
	}

	return nil
}
