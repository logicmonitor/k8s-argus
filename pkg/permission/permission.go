package permission

import (
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	enable  = true
	disable = false

	client kubernetes.Interface

	deploymentPermissionFlag              *bool
	horizontalPodAutoscalerPermissionFlag *bool
	mu                                    sync.Mutex
)

// Init is a function than init the permission context
func Init(k8sClient kubernetes.Interface) {
	mu.Lock()
	defer mu.Unlock()
	client = k8sClient
}

// HasDeploymentPermissions is a function that check if the deployment resource has permissions
// nolint: dupl
func HasDeploymentPermissions() bool {
	if deploymentPermissionFlag != nil {
		return *deploymentPermissionFlag
	}
	mu.Lock()
	defer mu.Unlock()
	_, err := client.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		deploymentPermissionFlag = &disable
		logrus.Errorf("Failed to list deployments: %+v", err)
	} else {
		deploymentPermissionFlag = &enable
	}

	return *deploymentPermissionFlag
}

// HasHorizontalPodAutoscalerPermissions is a function that checks if the Horizontal Pod Autoscaler resource has permissions
// nolint: dupl
func HasHorizontalPodAutoscalerPermissions() bool {
	if horizontalPodAutoscalerPermissionFlag != nil {
		return *horizontalPodAutoscalerPermissionFlag
	}

	mu.Lock()
	defer mu.Unlock()

	_, err := client.AutoscalingV1().HorizontalPodAutoscalers(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		horizontalPodAutoscalerPermissionFlag = &disable

		logrus.Errorf("Failed to list horizontalPodAutoscalers: %+v", err)
	} else {
		horizontalPodAutoscalerPermissionFlag = &enable
	}

	return *horizontalPodAutoscalerPermissionFlag
}

// HasPermissions has permission
func HasPermissions(rt enums.ResourceType) bool {
	switch rt {
	case enums.Deployments:

		return HasDeploymentPermissions()
	case enums.Hpas:

		return HasHorizontalPodAutoscalerPermissions()
	case enums.ETCD, enums.Namespaces, enums.Nodes, enums.Pods, enums.Services, enums.Unknown:

		return true
	default:

		return true
	}
}
