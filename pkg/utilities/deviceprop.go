package utilities

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/lm-sdk-go/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetResourcePropertyValue get resource property value by property name
func GetResourcePropertyValue(resource *models.Device, propertyName string) string {
	if resource == nil {
		return ""
	}
	if len(resource.CustomProperties) > 0 {
		for _, cp := range resource.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}
	if len(resource.SystemProperties) > 0 {
		for _, cp := range resource.SystemProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}

	return ""
}

// GetResourceGroupPropertyValue get resource property value by property name
func GetResourceGroupPropertyValue(resource *models.DeviceGroup, propertyName string) string {
	if resource == nil {
		return ""
	}
	if len(resource.CustomProperties) > 0 {
		for _, cp := range resource.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}

	return ""
}

// TrimName it will trim the name to 244 char if greater than 244
func TrimName(name string) string {
	if len(name) > constants.MaxResourceLength {
		name = name[:constants.MaxResourceLength]
	}

	return name
}

// GetResourceType get rt
func GetResourceType(resource *models.Device) (enums.ResourceType, error) {
	categories := GetResourcePropertyValue(resource, constants.K8sSystemCategoriesPropertyKey)
	list := strings.Split(categories, ",")
	for _, category := range list {
		if strings.HasPrefix(category, "Kubernetes") {
			category = strings.TrimPrefix(category, "Kubernetes")
			category = strings.TrimSuffix(category, "Conflict")
			category = strings.TrimSuffix(category, "Deleted")
			return enums.ParseResourceType(category)
		}
	}

	return enums.Unknown, fmt.Errorf("no valid category found in system.categories")
}

func IsArgusPodObject(lctx *lmctx.LMContext, rt enums.ResourceType, meta *metav1.PartialObjectMetadata) bool {
	if rt == enums.Pods {
		for k, v := range meta.Labels {
			if k == "app" && (v == constants.Argus || v == constants.CollectorsetController) {
				return true
			}
		}
	}
	return false
}
