package utilities

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
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

//IsConflictingDevice checks wheather there is conflicts in device names.
func IsConflictingDevice(device *models.Device) bool {
	nameconflictProp := GetPropertyValue(device, constants.K8sDeviceNameConflictPropertyKey)
	result, err := strconv.ParseBool(nameconflictProp)
	if err != nil {
		return false
	}
	return result
}

// GetDesiredDisplayNameByResourceAndConfig returns desired display name based on FullDisplayNameIncludeClusterName and FullDisplayNameIncludeNamespace properties.
func GetDesiredDisplayNameByResourceAndConfig(name, namespace, clusterName, resource string, displayNameIncludeNamespace, displayNameIncludeClusterName bool) string {
	desiredName := getNameWithResourceType(name, resource)
	if displayNameIncludeClusterName {
		if strings.EqualFold(resource, constants.Nodes) {
			return fmt.Sprintf("%s-%s", desiredName, clusterName)
		}
		return fmt.Sprintf("%s-%s-%s", desiredName, namespace, clusterName)
	}
	if displayNameIncludeNamespace && !strings.EqualFold(resource, constants.Nodes) {
		return fmt.Sprintf("%s-%s", desiredName, namespace)
	}
	return desiredName
}

// GetFullDisplayName returns complete display name for a device.
func GetFullDisplayName(device *models.Device, resource, clusterName string) string {
	displayNameWithNamespace := GetDisplayNameWithNamespace(device, resource)
	return fmt.Sprintf("%s-%s", displayNameWithNamespace, clusterName)
}

//GetDisplayNameWithNamespace return displayName in the format - name-type-namespace
func GetDisplayNameWithNamespace(device *models.Device, resource string) string {
	nameWithResourceType := getNameWithResourceType(GetPropertyValue(device, constants.K8sDeviceNamePropertyKey), resource)
	namespace := GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey)
	if strings.EqualFold(resource, constants.Nodes) {
		return nameWithResourceType
	}
	displayName := fmt.Sprintf("%s-%s", nameWithResourceType, namespace)
	return displayName
}

//GetNameWithResourceType return resourcename with its respetive type.
func getNameWithResourceType(name, resource string) string {
	switch strings.ToLower(resource) {
	case constants.Pods:
		return fmt.Sprintf("%s-%s", name, "pod")
	case constants.Deployments:
		return fmt.Sprintf("%s-%s", name, "deploy")
	case constants.Services:
		return fmt.Sprintf("%s-%s", name, "svc")
	case constants.Nodes:
		return fmt.Sprintf("%s-%s", name, "node")
	case constants.HorizontalPodAutoScalers:
		return fmt.Sprintf("%s-%s", name, "hpa")
	}
	return name
}
