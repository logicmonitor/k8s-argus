package utilities

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/lm-sdk-go/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// selfLink utility to create self links
// Haven't export as there is no usecase/caller with these params
func selfLink(namespaced bool, apiVersion string, kind string, namespace string, name string) string {
	var selfLinkAPIPrefix string
	if apiVersion == constants.K8sAPIVersionV1 {
		selfLinkAPIPrefix = "/api"
	} else {
		selfLinkAPIPrefix = "/apis"
	}
	if namespaced {
		if namespace == "" {
			return ""
		}

		return fmt.Sprintf(selfLinkAPIPrefix+"/%s/namespaces/%s/%s/%s", apiVersion, namespace, kind, name)
	}

	return fmt.Sprintf(selfLinkAPIPrefix+"/%s/%s/%s", apiVersion, kind, name)
}

// SelfLink utility to create self links
func SelfLink(namespaced bool, apiVersion string, kind string, objectMeta *metav1.PartialObjectMetadata) string {
	return selfLink(namespaced, apiVersion, kind, objectMeta.Namespace, objectMeta.Name)
}

func ResourceCacheContainerValue(resource *models.Device) string {
	rt, err := GetResourceType(resource)
	if err != nil {
		return ""
	}
	if rt.IsNamespaceScopedResource() {
		return GetResourcePropertyValue(resource, constants.K8sResourceNamespacePropertyKey)
	}
	return constants.ClusterScopedGroupName
}

func ResourceContainerValueFromMeta(rt enums.ResourceType, metaObject *metav1.PartialObjectMetadata) string {
	if rt.IsNamespaceScopedResource() {
		return metaObject.Namespace
	}
	return constants.ClusterScopedGroupName
}
