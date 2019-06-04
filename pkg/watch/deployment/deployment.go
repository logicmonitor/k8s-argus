// Package deployment provides the logic for mapping a Kubernetes deployment to a
// LogicMonitor w.
package deployment

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta2"
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

// ApiVersion is a function that implements the Watcher interface.
func (w *Watcher) ApiVersion() string {
	return constants.K8sApiVersion_apps_v1beta2
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &v1beta2.Deployment{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		deployment := obj.(*v1beta2.Deployment)
		log.Infof("Handling add deployment event: %s", deployment.Name)
		w.add(deployment)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1beta2.Deployment)
		new := newObj.(*v1beta2.Deployment)

		w.update(old, new)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		deployment := obj.(*v1beta2.Deployment)
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
func (w *Watcher) add(deployment *v1beta2.Deployment) {
	if _, err := w.Add(
		w.args(deployment, constants.DeploymentCategory)...,
	); err != nil {
		log.Errorf("Failed to add deployment %q: %v", fmtDeploymentDisplayName(deployment), err)
		return
	}
	log.Infof("Added deployment %q", fmtDeploymentDisplayName(deployment))
}

func (w *Watcher) update(old, new *v1beta2.Deployment) {
	if _, err := w.UpdateAndReplaceByDisplayName(
		fmtDeploymentDisplayName(old),
		w.args(new, constants.DeploymentCategory)...,
	); err != nil {
		log.Errorf("Failed to update deployment %q: %v", fmtDeploymentDisplayName(new), err)
		return
	}
	log.Infof("Updated deployment %q", fmtDeploymentDisplayName(old))
}

func (w *Watcher) move(deployment *v1beta2.Deployment) {
	if _, err := w.UpdateAndReplaceFieldByDisplayName(fmtDeploymentDisplayName(deployment), constants.CustomPropertiesFieldName, w.args(deployment, constants.DeploymentDeletedCategory)...); err != nil {
		log.Errorf("Failed to move deployment %q: %v", fmtDeploymentDisplayName(deployment), err)
		return
	}
	log.Infof("Moved deployment %q", fmtDeploymentDisplayName(deployment))
}

func (w *Watcher) args(deployment *v1beta2.Deployment, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, deployment.Labels)
	return []types.DeviceOption{
		w.Name(deployment.Name),
		w.ResourceLabels(deployment.Labels),
		w.DisplayName(fmtDeploymentDisplayName(deployment)),
		w.SystemCategories(categories),
		w.Auto("name", deployment.Name),
		w.Auto("namespace", deployment.Namespace),
		w.Auto("selflink", deployment.SelfLink),
		w.Auto("uid", string(deployment.UID)),
	}
}

// fmtDeploymentDisplayName implements the conversion for the deployment display name
func fmtDeploymentDisplayName(deployment *v1beta2.Deployment) string {
	return fmt.Sprintf("%s.%s.deploy-%s", deployment.Name, deployment.Namespace, string(deployment.UID))
}

// GetDeploymentsMap implements the getting deployments map info from k8s
func GetDeploymentsMap(k8sClient *kubernetes.Clientset, namespace string) (map[string]string, error) {
	deploymentsMap := make(map[string]string)
	deploymentList, err := k8sClient.AppsV1beta2().Deployments(namespace).List(v1.ListOptions{})
	if err != nil || deploymentList == nil {
		log.Warnf("Failed to get the deployments from k8s")
		return nil, err
	}
	for _, deploymentInfo := range deploymentList.Items {
		deploymentsMap[fmtDeploymentDisplayName(&deploymentInfo)] = deploymentInfo.Name
	}

	return deploymentsMap, nil
}
