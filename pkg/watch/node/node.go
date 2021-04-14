// Package node provides the logic for mapping a Kubernetes Node to a
// LogicMonitor device.
package node

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	resource = "nodes"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct {
	types.DeviceManager
	DeviceGroups map[string]int32
	LMClient     *client.LMSdkGo
	*types.WConfig
}

// APIVersion is a function that implements the Watcher interface.
func (w *Watcher) APIVersion() string {
	return constants.K8sAPIVersionV1
}

// Enabled is a function that check the resource can watch.
func (w *Watcher) Enabled() bool {
	return true
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// Namespaced returns true if resource is namespaced
func (w *Watcher) Namespaced() bool {
	return false
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &v1.Node{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + node.Name}))
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)

		log.Debugf("Handling add node event: %s", w.getDesiredDisplayName(node))

		// Require an IP address.
		if getInternalAddress(node.Status.Addresses) == nil {
			return
		}
		w.add(lctx, node)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Node)
		new := newObj.(*v1.Node)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + w.getDesiredDisplayName(old)}))
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)

		log.Debugf("Handling update node event: %s", w.getDesiredDisplayName(old))

		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := getInternalAddress(old.Status.Addresses)
		newInternalAddress := getInternalAddress(new.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			w.add(lctx, new)
			return
		}
		// Covers the case when the old node is in the process of terminating
		// and the new node is coming up to replace it.
		// if oldInternalAddress.Address != newInternalAddress.Address {
		w.update(lctx, old, new)
		// }
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + w.getDesiredDisplayName(node)}))
		log := lmlog.Logger(lctx)

		log.Debugf("Handling delete node event: %s", w.getDesiredDisplayName(node))

		// nolint: dupl
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(node),
				fmtNodeDisplayName(node, w.Config().ClusterName)); err != nil {
				log.Errorf("Failed to delete node: %v", err)
				return
			}
			log.Infof("Deleted node %s", w.getDesiredDisplayName(node))
			return
		}

		// Move the node.
		w.move(lctx, node)
	}
}

// nolint: dupl
func (w *Watcher) add(lctx *lmctx.LMContext, node *v1.Node) {
	log := lmlog.Logger(lctx)
	n, err := w.Add(lctx, w.Resource(), node.Labels, w.args(node, constants.NodeCategory)...)
	if err != nil {
		log.Errorf("Failed to add node %q: %v", w.getDesiredDisplayName(node), err)
		return
	}
	if n == nil {
		log.Debugf("node %q is not added as it is mentioned for filtering.", w.getDesiredDisplayName(node))
		return
	}

	log.Infof("Added node %q", *n.DisplayName)
	w.createRoleDeviceGroup(lctx, node.Labels)
}

func (w *Watcher) nodeUpdateFilter(old, new *v1.Node) types.UpdateFilter {
	return func() bool {
		return getInternalAddress(old.Status.Addresses) != getInternalAddress(new.Status.Addresses)
	}
}

func (w *Watcher) update(lctx *lmctx.LMContext, old, new *v1.Node) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(old),
		fmtNodeDisplayName(old, w.Config().ClusterName), w.nodeUpdateFilter(old, new),
		new.Labels, w.args(new, constants.NodeCategory)...); err != nil {
		log.Errorf("Failed to update node %q: %v", w.getDesiredDisplayName(new), err)
	} else {
		log.Infof("Updated node %q", w.getDesiredDisplayName(old))
	}

	// determine if we need to add a new node role device group
	oldLabel, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, old.Labels)
	newLabel, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, new.Labels)
	if oldLabel != newLabel {
		w.createRoleDeviceGroup(lctx, new.Labels)
	}
}

// nolint: dupl
func (w *Watcher) move(lctx *lmctx.LMContext, node *v1.Node) {
	log := lmlog.Logger(lctx)
	if _, err := w.MoveToDeletedGroup(lctx, w.Resource(), w.getDesiredDisplayName(node),
		fmtNodeDisplayName(node, w.Config().ClusterName), node.DeletionTimestamp, w.args(node, constants.NodeDeletedCategory)...); err != nil {
		log.Errorf("Failed to move node %q: %v", w.getDesiredDisplayName(node), err)
		return
	}
	log.Infof("Moved node %q", w.getDesiredDisplayName(node))
}

func (w *Watcher) args(node *v1.Node, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(getInternalAddress(node.Status.Addresses).Address),
		w.ResourceLabels(node.Labels),
		w.DisplayName(w.getDesiredDisplayName(node)),
		w.SystemCategories(category),
		w.Auto("name", node.Name),
		w.Auto("selflink", util.SelfLink(w.Namespaced(), w.APIVersion(), w.Resource(), node.ObjectMeta)),
		w.Auto("uid", string(node.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(node.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, node.Name),
	}
}

func fmtNodeDisplayName(node *v1.Node, clusterName string) string {
	return fmt.Sprintf("%s-node-%s", node.Name, clusterName)
}

func (w *Watcher) getDesiredDisplayName(node *v1.Node) string {
	return w.DeviceManager.GetDesiredDisplayName(node.Name, node.Namespace, constants.Nodes)
}

// getInternalAddress finds the node's internal address.
func getInternalAddress(addresses []v1.NodeAddress) *v1.NodeAddress {
	var hostname *v1.NodeAddress
	for i := range addresses {
		address := addresses[i]
		if address.Type == v1.NodeInternalIP {
			return &address
		}
		if address.Type == v1.NodeHostName {
			hostname = &address
		}
	}
	//if there is no internal IP for this node, the host name will be used
	return hostname
}

func (w *Watcher) createRoleDeviceGroup(lctx *lmctx.LMContext, labels map[string]string) {
	log := lmlog.Logger(lctx)
	label, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, labels)
	if label == "" {
		return
	}
	role := strings.Replace(label, constants.LabelNodeRole, "", -1)

	if devicegroup.Exists(w.DeviceGroups[constants.ClusterDeviceGroupPrefix+w.Config().ClusterName], role, w.LMClient) {
		log.Infof("Device group for node role %q already exists", role)
		return
	}

	opts := &devicegroup.Options{
		ParentID:              w.DeviceGroups[constants.NodeDeviceGroupName],
		Name:                  role,
		DisableAlerting:       w.Config().DisableAlerting,
		AppliesTo:             devicegroup.NewAppliesToBuilder().Exists(constants.LabelCustomPropertyPrefix + label).And().HasCategory(constants.NodeCategory).And().Auto("clustername").Equals(w.Config().ClusterName),
		Client:                w.LMClient,
		DeleteDevices:         w.Config().DeleteDevices,
		AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().Exists(constants.LabelCustomPropertyPrefix + label).And().HasCategory(constants.NodeDeletedCategory).And().Auto("clustername").Equals(w.Config().ClusterName),
		CustomProperties:      devicegroup.NewPropertyBuilder(),
	}

	log.Debugf("%v", opts)

	_, err := devicegroup.Create(opts)
	if err != nil {
		log.Errorf("Failed to add device group for node role to %q: %v", role, err)
		return
	}

	log.Infof("Added device group for node role %q", role)
}

// GetNodesMap implements the getting nodes map info from k8s
func GetNodesMap(k8sClient kubernetes.Interface, clusterName string) (map[string]string, error) {
	nodesMap := make(map[string]string)
	nodeList, err := k8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil || nodeList == nil {
		return nil, err
	}
	for _, nodeInfo := range nodeList.Items {
		address := getInternalAddress(nodeInfo.Status.Addresses)
		if address == nil {
			continue
		}
		nodesMap[fmtNodeDisplayName(&nodeInfo, clusterName)] = address.Address
	}

	return nodesMap, nil
}
