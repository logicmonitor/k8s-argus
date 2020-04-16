// Package deployment provides the logic for mapping a Kubernetes deployment to a
// LogicMonitor w.
package deployment

import (
	"fmt"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	resource = "deployments"
)

// Watcher represents a watcher type that watches deployments.
type Watcher struct {
	types.DeviceManager
}

// APIVersion is a function that implements the Watcher interface.
func (w *Watcher) APIVersion() string {
	return constants.K8sAPIVersionAppsV1
}

// Enabled is a function that check the resource can watch.
func (w *Watcher) Enabled() bool {
	return permission.HasDeploymentPermissions()
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &appsv1.Deployment{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		deployment := obj.(*appsv1.Deployment)
		log.Infof("Handling add deployment event: %s", deployment.Name)
		w.add(deployment)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*appsv1.Deployment)
		new := newObj.(*appsv1.Deployment)

		w.update(old, new)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		deployment := obj.(*appsv1.Deployment)
		log.Debugf("Handling delete deployment event: %s", deployment.Name)
		// Delete the deployment.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(fmtDeploymentDisplayName(deployment)); err != nil {
				log.Errorf("Failed to delete deployment: %v", err)
				return
			}
			log.Infof("Deleted deployment %s", deployment.Name)
			return
		}

		// Move the deployment.
		w.move(deployment)
	}
}

// nolint: dupl
func (w *Watcher) add(deployment *appsv1.Deployment) {
	if _, err := w.Add(
		w.args(deployment, constants.DeploymentCategory)...,
	); err != nil {
		log.Errorf("Failed to add deployment %q: %v", fmtDeploymentDisplayName(deployment), err)
		return
	}
	log.Infof("Added deployment %q", fmtDeploymentDisplayName(deployment))
}

func (w *Watcher) update(old, new *appsv1.Deployment) {
	if _, err := w.UpdateAndReplaceByDisplayName(
		fmtDeploymentDisplayName(old),
		w.args(new, constants.DeploymentCategory)...,
	); err != nil {
		log.Errorf("Failed to update deployment %q: %v", fmtDeploymentDisplayName(new), err)
		return
	}
	log.Infof("Updated deployment %q", fmtDeploymentDisplayName(old))
}

func (w *Watcher) move(deployment *appsv1.Deployment) {
	if _, err := w.UpdateAndReplaceFieldByDisplayName(fmtDeploymentDisplayName(deployment), constants.CustomPropertiesFieldName, w.args(deployment, constants.DeploymentDeletedCategory)...); err != nil {
		log.Errorf("Failed to move deployment %q: %v", fmtDeploymentDisplayName(deployment), err)
		return
	}
	log.Infof("Moved deployment %q", fmtDeploymentDisplayName(deployment))
}

func (w *Watcher) args(deployment *appsv1.Deployment, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(deployment.Name),
		w.ResourceLabels(deployment.Labels),
		w.DisplayName(fmtDeploymentDisplayName(deployment)),
		w.SystemCategories(category),
		w.Auto("name", deployment.Name),
		w.Auto("namespace", deployment.Namespace),
		w.Auto("selflink", deployment.SelfLink),
		w.Auto("uid", string(deployment.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(deployment.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, fmtDeploymentDisplayName(deployment)),
	}
}

// fmtDeploymentDisplayName implements the conversion for the deployment display name
func fmtDeploymentDisplayName(deployment *appsv1.Deployment) string {
	return fmt.Sprintf("%s.%s.deploy-%s", deployment.Name, deployment.Namespace, string(deployment.UID))
}

// GetDeploymentsMap implements the getting deployments map info from k8s
func GetDeploymentsMap(k8sClient *kubernetes.Clientset, namespace string) (map[string]string, error) {
	deploymentsMap := make(map[string]string)
	deploymentList, err := k8sClient.AppsV1().Deployments(namespace).List(v1.ListOptions{})
	if err != nil || deploymentList == nil {
		log.Warnf("Failed to get the deployments from k8s")
		return nil, err
	}
	for _, deploymentInfo := range deploymentList.Items {
		deploymentsMap[fmtDeploymentDisplayName(&deploymentInfo)] = deploymentInfo.Name
	}

	return deploymentsMap, nil
}
