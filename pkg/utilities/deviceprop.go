package utilities

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// GetPropertyValue get device property value by property name
func GetPropertyValue(device *models.Device, propertyName string) string {
	if device == nil {
		return ""
	}
	if len(device.CustomProperties) > 0 {
		for _, cp := range device.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}
	if len(device.SystemProperties) > 0 {
		for _, cp := range device.SystemProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}

	return ""
}

// IsConflictingDevice checks wheather there is conflicts in device names.
func IsConflictingDevice(device *models.Device, resourceType enums.ResourceType) bool {
	sysCategory := GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)

	return strings.Contains(sysCategory, GetConflictCategoryByResourceType(resourceType))
}

// GetDisplayName returns desired display name based on FullDisplayNameIncludeClusterName and FullDisplayNameIncludeNamespace properties.
func GetDisplayName(name, namespace string, resource enums.ResourceType, config *config.Config) string {
	clusterName := config.ClusterName
	includeNamespace := config.FullDisplayNameIncludeNamespace
	includeClusterName := config.FullDisplayNameIncludeClusterName

	desiredName := getNameWithResourceType(name, resource)

	if includeClusterName {
		if resource.IsNamespaceScopedResource() {
			return fmt.Sprintf("%s-%s-%s", desiredName, namespace, clusterName)
		}

		return fmt.Sprintf("%s-%s", desiredName, clusterName)
	}
	if includeNamespace && resource.IsNamespaceScopedResource() {
		return fmt.Sprintf("%s-%s", desiredName, namespace)
	}

	return desiredName
}

// GetFullDisplayName returns complete display name for a device.
func GetFullDisplayName(device *models.Device, resource enums.ResourceType, clusterName string) string {
	displayNameWithNamespace := GetDisplayNameWithNamespace(device, resource)

	return fmt.Sprintf("%s-%s", displayNameWithNamespace, clusterName)
}

// GetDisplayNameWithNamespace returns displayName in the format - name-type-namespace
func GetDisplayNameWithNamespace(device *models.Device, resource enums.ResourceType) string {
	nameWithResourceType := getNameWithResourceType(GetPropertyValue(device, constants.K8sDeviceNamePropertyKey), resource)
	namespace := GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey)
	if !resource.IsNamespaceScopedResource() {
		return nameWithResourceType
	}
	displayName := fmt.Sprintf("%s-%s", nameWithResourceType, namespace)

	return displayName
}

// GetNameWithResourceType returns resourcename with its respetive type.
func getNameWithResourceType(name string, resource enums.ResourceType) string {
	resourceType := enums.ShortResourceType(resource)

	return fmt.Sprintf("%s-%s", name, resourceType.String())
}

// GetConflictCategoryByResourceType returns conflict system category by its respetive type.
func GetConflictCategoryByResourceType(resource enums.ResourceType) string {
	return resource.GetConflictsCategory()
}

// TrimName it will trim the name to 244 char if greater than 244
func TrimName(name string) string {
	if len(name) > constants.MaxResourceLength {
		name = name[:constants.MaxResourceLength]
	}

	return name
}

// GetNameWithResourceTypeAndNamespace returns name with resource_type and namespace
func GetNameWithResourceTypeAndNamespace(name string, resource enums.ResourceType, namespace string) string {
	resourceType := enums.ShortResourceType(resource)
	if resource.IsNamespaceScopedResource() {
		return fmt.Sprintf("%s-%s-%s", name, resourceType.String(), namespace)
	}

	return fmt.Sprintf("%s-%s", name, resourceType.String())
}

// GetResourceType get rt
func GetResourceType(device *models.Device) (enums.ResourceType, error) {
	categories := GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)
	var l enums.ResourceType
	if strings.Contains(categories, constants.PodCategory) {
		l = enums.Pods
	}
	if strings.Contains(categories, constants.DeploymentCategory) {
		l = enums.Deployments
	}
	if strings.Contains(categories, constants.ServiceCategory) {
		l = enums.Services
	}
	if strings.Contains(categories, constants.NodeCategory) {
		l = enums.Nodes
	}
	if strings.Contains(categories, constants.HorizontalPodAutoscalerCategory) {
		l = enums.Hpas
	}
	if l != enums.Unknown {
		return l, nil
	}

	return l, fmt.Errorf("no valid category found in system.categories")
}
