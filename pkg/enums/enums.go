package enums

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceType resource
type ResourceType uint32

// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
// START This comment starts the section to be changed when enum added or modified
const (
	Unknown ResourceType = iota
	Pods
	Deployments
	Services
	Hpas
	Nodes
	ETCD
	Namespaces
)

// ALLResourceTypes All resource type slice
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
var ALLResourceTypes = []ResourceType{
	Pods,
	Deployments,
	Services,
	Hpas,
	Nodes,
	ETCD,
	Namespaces,
}

// MarshalText marshals
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType ResourceType) MarshalText() ([]byte, error) {
	switch resourceType {
	case Pods:

		return []byte("pods"), nil
	case Deployments:

		return []byte("deployments"), nil
	case Services:

		return []byte("services"), nil
	case Hpas:

		return []byte("horizontalpodautoscalers"), nil
	case Nodes:

		return []byte("nodes"), nil
	case ETCD:

		return []byte("etcd"), nil
	case Namespaces:

		return []byte("namespaces"), nil
	case Unknown:

		return []byte("unknown"), nil
	}

	// do not put call for String method again here, leads to loop and stackoverflow
	// gracefully converting it to uint32. otherwise goes in loop and leads to stackoverflow

	return nil, fmt.Errorf("not a valid ResourceType to marshal %v", uint32(resourceType))
}

// ParseResourceType takes a string level and returns the Logrus log level constant.
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func ParseResourceType(resourceType string) (ResourceType, error) {
	switch strings.ToLower(resourceType) {
	case "pods", "pod", "po":

		return Pods, nil
	case "deployments", "deployment", "deploy":

		return Deployments, nil
	case "services", "service", "svc":

		return Services, nil
	case "hpas", "horizontalpodautoscalers", "horizontalpodautoscaler", "hpa":

		return Hpas, nil
	case "nodes", "node":

		return Nodes, nil
	case "etcd":

		return ETCD, nil
	case "namespaces", "namespace", "ns":

		return Namespaces, nil
	}

	return Unknown, fmt.Errorf("not a valid ResourceType to parse: %s", resourceType)
}

// MarshalText marshals
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType ShortResourceType) MarshalText() ([]byte, error) {
	switch ResourceType(resourceType) {
	case Pods:

		return []byte("pod"), nil
	case Deployments:

		return []byte("deploy"), nil
	case Services:

		return []byte("svc"), nil
	case Hpas:

		return []byte("hpa"), nil
	case Nodes:

		return []byte("node"), nil
	case ETCD:

		return []byte("etcd"), nil
	case Namespaces:

		return []byte("ns"), nil
	case Unknown:

		return []byte("unknown"), nil
	}

	return nil, fmt.Errorf("not a valid ShortResourceType to marshal %d", resourceType)
}

// ParseShortResourceType takes a string level and returns the Logrus log level constant.
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func ParseShortResourceType(shortResourceType string) (ShortResourceType, error) {
	var l ResourceType
	switch strings.ToLower(shortResourceType) {
	case "pod":
		l = Pods
	case "deploy":
		l = Deployments
	case "svc":
		l = Services
	case "hpa", "horizontalpodautoscaler":
		l = Hpas
	case "node":
		l = Nodes
	case "etcd":
		l = ETCD
	case "ns":
		l = Namespaces
	default:

		return ShortResourceType(Unknown), fmt.Errorf("not a valid ShortResourceType to parse: %q", shortResourceType)
	}

	return ShortResourceType(l), nil
}

// K8SObjectType returns runtime.Object to create watcher
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) K8SObjectType() runtime.Object {
	switch *resourceType {
	case Pods:

		return &corev1.Pod{} // nolint: exhaustivestruct
	case Deployments:

		return &appsv1.Deployment{} // nolint: exhaustivestruct
	case Services:

		return &corev1.Service{} // nolint: exhaustivestruct
	case Hpas:

		return &autoscalingv1.HorizontalPodAutoscaler{} // nolint: exhaustivestruct
	case Nodes:

		return &corev1.Node{} // nolint: exhaustivestruct
	case Namespaces:

		return &corev1.Namespace{} // nolint: exhaustivestruct
	case ETCD, Unknown:

		return nil
	default:

		return nil
	}
}

// K8SAPIVersion returns runtime.Object to create watcher
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) K8SAPIVersion() string {
	switch *resourceType {
	// core v1 version
	case Pods, Services, Nodes, Namespaces:

		return constants.K8sAPIVersionV1
	// apps api group apps/v1 version
	case Deployments:

		return constants.K8sAPIVersionAppsV1
	// autoscaling api group v1 version
	case Hpas:

		return constants.K8sAutoscalingV1
	case ETCD, Unknown:

		return ""
	default:

		return ""
	}
}

// ObjectMeta returns object meta from interface object
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) ObjectMeta(obj interface{}) *metav1.ObjectMeta {
	switch *resourceType {
	case Pods:

		return &obj.(*corev1.Pod).ObjectMeta
	case Deployments:

		return &obj.(*appsv1.Deployment).ObjectMeta
	case Services:

		return &obj.(*corev1.Service).ObjectMeta
	case Hpas:

		return &obj.(*autoscalingv1.HorizontalPodAutoscaler).ObjectMeta
	case Nodes:

		return &obj.(*corev1.Node).ObjectMeta
	case Namespaces:

		return &obj.(*corev1.Namespace).ObjectMeta
	case ETCD, Unknown:

		return nil
	default:

		return nil
	}
}

// IsNamespaceScopedResource returns true if namespace scoped
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) IsNamespaceScopedResource() bool {
	switch *resourceType {
	case Pods, Deployments, Services, Hpas, Namespaces:

		return true
	case Nodes, ETCD, Unknown:

		return false
	default:
		// TODO: whether defaults to true or false?

		return false
	}
}

// GetCategory returns category name for conflicts group
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) GetCategory() string {
	switch *resourceType {
	case Pods:

		return constants.PodCategory
	case Deployments:

		return constants.DeploymentCategory
	case Services:

		return constants.ServiceCategory
	case Nodes:

		return constants.NodeCategory
	case Hpas:

		return constants.HorizontalPodAutoscalerCategory
	case ETCD:

		return constants.EtcdCategory
	case Namespaces:

		return constants.NamespaceCategory
	case Unknown:
		unknown := Unknown

		return unknown.String()
	default:
		unknown := Unknown

		return unknown.String()
	}
}

// APIGroup returns string name
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) APIGroup() string {
	switch *resourceType {
	// core v1 version
	case Pods, Services, Nodes, Namespaces:

		return ""
	// apps api group apps/v1 version
	case Deployments:

		return "apps"
	// autoscaling api group v1 version
	case Hpas:

		return "autoscaling"
	case ETCD, Unknown:

		return ""
	default:

		return ""
	}
}

// IsK8SPingResource returns true if resource can be pinged using system.ips prop - generic host status uses hostname, for special handling k8s ping ds uses system.ips
// NOTE: RESOURCE_MODIFICATION need to change when adding/deleting resource for monitoring
func (resourceType *ResourceType) IsK8SPingResource() bool {
	switch *resourceType {
	case Pods, Unknown:

		return true
	case Deployments, Services, Hpas, Nodes, ETCD, Namespaces:

		return false
	default:

		return true
	}
}

// END sections ends here when enum added or modified
//
//
//
//
//
//
//
//
//

// String returns string name
func (resourceType ResourceType) String() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return str
	}

	return ""
}

// FQName returns string name
func (resourceType *ResourceType) FQName(name string) string {
	if apiGroup := resourceType.APIGroup(); apiGroup != "" {
		return fmt.Sprintf("%s.%s/%s", resourceType.Singular(), apiGroup, name)
	}

	return fmt.Sprintf("%s/%s", resourceType.Singular(), name)
}

// Singular returns string name
func (resourceType *ResourceType) Singular() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return strings.TrimSuffix(str, "s")
	}

	return ""
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (resourceType *ResourceType) UnmarshalText(text []byte) error {
	l, err := ParseResourceType(string(text))
	if err != nil {
		return err
	}

	*resourceType = l

	return nil
}

// LMName lmname
func (resourceType ResourceType) LMName(meta *metav1.ObjectMeta) string {
	s := ShortResourceType(resourceType)
	if resourceType.IsNamespaceScopedResource() {
		return fmt.Sprintf("%s-%s-%s", meta.Name, s.String(), meta.Namespace)
	}

	return fmt.Sprintf("%s-%s", meta.Name, s.String())
}

// ShortResourceType to specifically use as short resource type
type ShortResourceType ResourceType

// UnmarshalText implements encoding.TextUnmarshaler.
func (resourceType *ShortResourceType) UnmarshalText(text []byte) error {
	l, err := ParseShortResourceType(string(text))
	if err != nil {
		return err
	}
	*resourceType = l

	return nil
}

// String returns string name
func (resourceType *ShortResourceType) String() string {
	bytes, err := json.Marshal(resourceType)
	if err != nil {
		return ""
	}
	if str, err := strconv.Unquote(fmt.Sprintf(`%s`, bytes)); err == nil {
		return str
	}

	return ""
}

// GetConflictsCategory returns category name for conflicts group
func (resourceType *ResourceType) GetConflictsCategory() string {
	return fmt.Sprintf("%s%s", resourceType.GetCategory(), "Conflict")
}

// GetDeletedCategory returns category name for conflicts group
func (resourceType *ResourceType) GetDeletedCategory() string {
	return fmt.Sprintf("%s%s", resourceType.GetCategory(), "Deleted")
}
