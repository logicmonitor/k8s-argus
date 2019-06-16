package permission

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	permissionFlagDefault  = 0
	permissionFlagEnable   = 1
	permissionFlagDisabled = -1
)

var (
	client *kubernetes.Clientset

	deploymentPermissionFlag = 0
)

// Init is a function than init the permission context
func Init(k8sClient *kubernetes.Clientset) {
	client = k8sClient
}

// HasDeploymentPermissions is a function that check if the deployment resource has permissions
func HasDeploymentPermissions() bool {
	if deploymentPermissionFlag != permissionFlagDefault {
		return deploymentPermissionFlag == permissionFlagEnable
	}
	_, err := client.AppsV1beta2().Deployments(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		deploymentPermissionFlag = permissionFlagDisabled
		log.Errorf("Failed to list deployments: %+v", err)
	} else {
		deploymentPermissionFlag = permissionFlagEnable
	}
	return deploymentPermissionFlag == permissionFlagEnable
}
