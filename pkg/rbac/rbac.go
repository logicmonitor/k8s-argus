package rbac

import (
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	rbacFlagDefault  = 0
	rbacFlagEnable   = 1
	rbacFlagDisabled = -1
)

var (
	client *kubernetes.Clientset

	deploymentRBACFlag = 0
)

// Init is a function than init the rbac context
func Init(k8sclient *kubernetes.Clientset) {
	client = k8sclient
}

// HasDeploymentRBAC is a function that check if the deployment resource has rbac permissions
func HasDeploymentRBAC() bool {
	if deploymentRBACFlag != rbacFlagDefault {
		return deploymentRBACFlag == rbacFlagEnable
	}
	_, err := client.AppsV1beta2().Deployments(v1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		deploymentRBACFlag = rbacFlagDisabled
		log.Errorf("Failed to list deployments: %+v", err)
	} else {
		deploymentRBACFlag = rbacFlagEnable
	}
	return deploymentRBACFlag == rbacFlagEnable
}
