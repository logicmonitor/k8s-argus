package k8s

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetAllK8SResources get all k8s resources present in cluster
func GetAllK8SResources() *resourcecache.Store {
	tmpStore := resourcecache.NewStore()
	conf, err := config.GetConfig()
	if err != nil {
		logrus.Errorf("Failed to read config")

		return nil
	}
	for _, rt := range enums.ALLResourceTypes {
		for _, meta := range GetAndStoreAll(rt) {
			displayName := util.GetDisplayNameNew(rt, &meta, conf) //nolint:gosec
			tmpStore.Set(types.ResourceName{
				Name:     meta.Name,
				Resource: rt,
			}, types.ResourceMeta{ // nolint: exhaustivestruct
				Container:   meta.Namespace,
				Labels:      meta.Labels,
				DisplayName: displayName,
				UID:         meta.UID,
			})
		}
	}

	return tmpStore
}

// GetAndStoreAll get
// nolint: gocyclo
func GetAndStoreAll(rt enums.ResourceType) []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)

	switch rt {
	case enums.Pods:
		result = getPodList()
	case enums.Deployments:
		result = getDeploymentList()
	case enums.Services:
		result = getServiceList()
	case enums.Nodes:
		result = getNodeList()
	case enums.Hpas:
		result = getHPAList()
	case enums.ETCD, enums.Namespaces, enums.Unknown:
		return result
	}

	return result
}

func getPodList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	podList, err := config.GetClientSet().CoreV1().Pods("").List(constants.DefaultListOptions)

	if err != nil || podList == nil {
		return result
	}
	for _, i := range podList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func getServiceList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	serviceList, err := config.GetClientSet().CoreV1().Services("").List(constants.DefaultListOptions)

	if err != nil || serviceList == nil {
		return result
	}
	for _, i := range serviceList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func getDeploymentList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	deploymentList, err := config.GetClientSet().AppsV1().Deployments("").List(constants.DefaultListOptions)

	if err != nil || deploymentList == nil {
		return result
	}
	for _, i := range deploymentList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func getNodeList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	nodeList, err := config.GetClientSet().CoreV1().Nodes().List(constants.DefaultListOptions)

	if err != nil || nodeList == nil {
		return result
	}
	for _, i := range nodeList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}

func getHPAList() []metav1.ObjectMeta {
	result := make([]metav1.ObjectMeta, 0)
	horizontalPodAutoscalerList, err := config.GetClientSet().AutoscalingV2beta2().HorizontalPodAutoscalers("").List(constants.DefaultListOptions)

	if err != nil || horizontalPodAutoscalerList == nil {
		return result
	}
	for _, i := range horizontalPodAutoscalerList.Items {
		result = append(result, i.ObjectMeta)
	}

	return result
}
