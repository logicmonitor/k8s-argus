package permission

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	enable  = true
	disable = false

	client kubernetes.Interface

	deploymentPermissionFlag              *bool
	horizontalPodAutoscalerPermissionFlag *bool
)

// Init is a function than init the permission context
func Init(k8sClient kubernetes.Interface) {
	client = k8sClient
}

// HasDeploymentPermissions is a function that check if the deployment resource has permissions
// nolint: dupl
func HasDeploymentPermissions() bool {
	if deploymentPermissionFlag != nil {
		return *deploymentPermissionFlag
	}
	_, err := client.AppsV1().Deployments(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		deploymentPermissionFlag = &disable
		log.Errorf("Failed to list deployments: %+v", err)
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
	_, err := client.AutoscalingV1().HorizontalPodAutoscalers(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		horizontalPodAutoscalerPermissionFlag = &disable
		log.Errorf("Failed to list horizontalPodAutoscalers: %+v", err)
	} else {
		horizontalPodAutoscalerPermissionFlag = &enable
	}
	return *horizontalPodAutoscalerPermissionFlag
}
