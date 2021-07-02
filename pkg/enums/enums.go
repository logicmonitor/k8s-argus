package enums

import (
	"fmt"
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

// TitlePlural returns string name in proper case
func (resourceType *ResourceType) TitlePlural() string {
	return fmt.Sprintf("%ss", resourceType.Title())
}

// Title returns string name in proper case
func (resourceType *ResourceType) Title() string {
	switch *resourceType {
	case Pods:

		return "Pod"
	case Deployments:

		return "Deployment"
	case Services:

		return "Service"
	case Hpas:

		return "HorizontalPodAutoscaler"
	case Nodes:

		return "Node"
	case ETCD:

		return "Etcd"
	case Namespaces:

		return "Namespace"
	case Unknown:

		return "Unknown"
	}

	return "Unknown"
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
	return fmt.Sprintf("%s%s", "Kubernetes", resourceType.Title())
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
