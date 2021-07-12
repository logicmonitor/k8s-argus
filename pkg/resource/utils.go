package resource

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func validateNewResource(lctx *lmctx.LMContext, resourceType enums.ResourceType, resource *models.Device, meta *metav1.PartialObjectMetadata) (bool, bool, error) {
	if resourceType.IsK8SPingResource() && util.GetResourcePropertyValue(resource, constants.K8sSystemIPsPropertyKey) == "" {
		return false, false, fmt.Errorf("property '%s' is empty for resource '%s'", constants.K8sSystemIPsPropertyKey, resourceType.FQName(meta.Name))
	}

	return false, true, nil
}
