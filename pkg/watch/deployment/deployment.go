// Package deployment provides the logic for mapping a Kubernetes deployment to a
// LogicMonitor w.
// nolint: dupl
package deployment

import (
	"fmt"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
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
	*types.WConfig
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
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + deployment.Name}))
		lctx = util.WatcherContext(lctx, w)
		log := lmlog.Logger(lctx)
		log.Infof("Handling add deployment event: %s", w.getDesiredDisplayName(deployment))
		w.add(lctx, deployment)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*appsv1.Deployment)
		new := newObj.(*appsv1.Deployment)

		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + old.Name}))
		lctx = util.WatcherContext(lctx, w)
		w.update(lctx, old, new)
	}
}

// DeleteFunc is a function that implements the Watcher interface.
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		deployment := obj.(*appsv1.Deployment)
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"device_id": resource + "-" + deployment.Name}))
		log := lmlog.Logger(lctx)
		log.Debugf("Handling delete deployment event: %s", deployment.Name)
		// Delete the deployment.
		if w.Config().DeleteDevices {
			if err := w.DeleteByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(deployment)); err != nil {
				log.Errorf("Failed to delete deployment: %v", err)
				return
			}
			log.Infof("Deleted deployment %s", deployment.Name)
			return
		}

		// Move the deployment.
		w.move(lctx, deployment)
	}
}

// nolint: dupl
func (w *Watcher) add(lctx *lmctx.LMContext, deployment *appsv1.Deployment) {
	log := lmlog.Logger(lctx)
	dep, err := w.Add(lctx, w.Resource(), deployment.Labels,
		w.args(deployment, constants.DeploymentCategory)...,
	)
	if err != nil {
		log.Errorf("Failed to add deployment %q: %v", w.getDesiredDisplayName(deployment), err)
		return
	}

	if dep == nil {
		log.Debugf("deployment %q is not added as it is mentioned for filtering.", w.getDesiredDisplayName(deployment))
		return
	}
	log.Infof("Added deployment %q", *dep.DisplayName)
}

func (w *Watcher) update(lctx *lmctx.LMContext, old, new *appsv1.Deployment) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceByDisplayName(lctx, w.Resource(),
		fmtDeploymentDisplayName(old, w.Config().ClusterName), nil, new.Labels,
		w.args(new, constants.DeploymentCategory)...,
	); err != nil {
		log.Errorf("Failed to update deployment %q: %v", w.getDesiredDisplayName(new), err)
		return
	}
	log.Infof("Updated deployment %q", w.getDesiredDisplayName(old))
}

// nolint: dupl
func (w *Watcher) move(lctx *lmctx.LMContext, deployment *appsv1.Deployment) {
	log := lmlog.Logger(lctx)
	if _, err := w.UpdateAndReplaceFieldByDisplayName(lctx, w.Resource(), w.getDesiredDisplayName(deployment), constants.CustomPropertiesFieldName, w.args(deployment, constants.DeploymentDeletedCategory)...); err != nil {
		log.Errorf("Failed to move deployment %q: %v", w.getDesiredDisplayName(deployment), err)
		return
	}
	log.Infof("Moved deployment %q", w.getDesiredDisplayName(deployment))
}

func (w *Watcher) args(deployment *appsv1.Deployment, category string) []types.DeviceOption {
	return []types.DeviceOption{
		w.Name(w.getDesiredDisplayName(deployment)),
		w.ResourceLabels(deployment.Labels),
		w.DisplayName(w.getDesiredDisplayName(deployment)),
		w.SystemCategories(category),
		w.Auto("name", deployment.Name),
		w.Auto("namespace", deployment.Namespace),
		w.Auto("selflink", deployment.SelfLink),
		w.Auto("uid", string(deployment.UID)),
		w.Custom(constants.K8sResourceCreatedOnPropertyKey, strconv.FormatInt(deployment.CreationTimestamp.Unix(), 10)),
		w.Custom(constants.K8sResourceNamePropertyKey, w.getDesiredDisplayName(deployment)),
	}
}

// fmtDeploymentDisplayName implements the conversion for the deployment display name
func fmtDeploymentDisplayName(deployment *appsv1.Deployment, clusterName string) string {
	return fmt.Sprintf("%s-deploy-%s-%s", deployment.Name, deployment.Namespace, clusterName)
}

func (w *Watcher) getDesiredDisplayName(deployment *appsv1.Deployment) string {
	return w.DeviceManager.GetDesiredDisplayName(deployment.Name, deployment.Namespace, constants.Deployments)
}

// GetDeploymentsMap implements the getting deployments map info from k8s
func GetDeploymentsMap(lctx *lmctx.LMContext, k8sClient kubernetes.Interface, namespace string, clusterName string) (map[string]string, error) {
	log := lmlog.Logger(lctx)
	deploymentsMap := make(map[string]string)
	deploymentList, err := k8sClient.AppsV1().Deployments(namespace).List(v1.ListOptions{})
	if err != nil || deploymentList == nil {
		log.Warnf("Failed to get the deployments from k8s")
		return nil, err
	}
	for i := range deploymentList.Items {
		deploymentsMap[fmtDeploymentDisplayName(&deploymentList.Items[i], clusterName)] = deploymentList.Items[i].Name
	}

	return deploymentsMap, nil
}
