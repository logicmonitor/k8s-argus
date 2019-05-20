package controller

import (
	"context"
	"fmt"
	"time"

	crv1alpha1 "github.com/logicmonitor/k8s-collectorset-controller/pkg/apis/v1alpha1"
	collectorsetclient "github.com/logicmonitor/k8s-collectorset-controller/pkg/client"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/config"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/distributor"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/distributor/roundrobin"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/policy"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/storage"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Controller is the Kubernetes controller object for LogicMonitor
// collectors.
type Controller struct {
	*collectorsetclient.Client
	CollectorSetScheme *runtime.Scheme
	LogicmonitorClient *lm.DefaultApi
	Storage            storage.Storage
}

// New instantiates and returns a Controller and an error if any.
func New(collectorsetconfig *config.Config, storage storage.Storage) (*Controller, error) {
	// Instantiate the Kubernetes in cluster config.
	restconfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Instantiate the CollectorSet client.
	client, collectorsetscheme, err := collectorsetclient.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// Instantiate the LogicMontitor API client.
	lmClient := newLMClient(collectorsetconfig.ID, collectorsetconfig.Key, collectorsetconfig.Account)

	// start a controller on instances of our custom resource
	c := &Controller{
		Client:             client,
		CollectorSetScheme: collectorsetscheme,
		LogicmonitorClient: lmClient,
		Storage:            storage,
	}

	return c, nil
}

// Run starts a CollectorSet resource controller.
func (c *Controller) Run(ctx context.Context) error {
	// Watch CollectorSet objects
	err := c.watch(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return ctx.Err()
}

func (c *Controller) watch(ctx context.Context) error {
	_, controller := cache.NewInformer(
		cache.NewListWatchFromClient(c.RESTClient, crv1alpha1.CollectorSetResourcePlural, apiv1.NamespaceAll, fields.Everything()),
		&crv1alpha1.CollectorSet{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.addFunc,
			UpdateFunc: c.updateFunc,
			DeleteFunc: c.deleteFunc,
		})

	go controller.Run(ctx.Done())

	return nil
}

func (c *Controller) addFunc(obj interface{}) {
	collectorset := obj.(*crv1alpha1.CollectorSet)
	log.Infof("Starting to create collectorset: %s", collectorset.Name)

	ids, err := CreateOrUpdateCollectorSet(collectorset, c.LogicmonitorClient, c.Clientset)
	if err != nil {
		log.Errorf("Failed to create collectorset: %v", err)
		return
	}
	log.Infof("CollectorSet %q has collectors %v", collectorset.Name, ids)

	log.Infof("Waiting for the collectors to register: %v", ids)
	if err = waitForCollectorsToRegister(c.LogicmonitorClient, ids); err != nil {
		log.Errorf("Failed to verify that collectors %v are registered: %v", ids, err)
		return
	}

	collectorsetCopy, err := c.updateCollectorSetStatus(collectorset, ids)
	if err != nil {
		log.Errorf("Failed to update collectorset status: %v", err)
		return
	}

	log.Infof("CollectorSet %q status is %q", collectorsetCopy.Name, collectorsetCopy.Status.State)

	if err = c.save(collectorsetCopy); err != nil {
		log.Errorf("Failed to save policy: %v", err)
	}

	log.Infof("Finished creating CollectorSet: %s", collectorsetCopy.Name)
}

// TODO: updating the collectorset ids in the add func will trigger this. We
// need to check for this case
func (c *Controller) updateFunc(oldObj, newObj interface{}) {
	_ = oldObj.(*crv1alpha1.CollectorSet)
	newcollectorset := newObj.(*crv1alpha1.CollectorSet)

	log.Infof("Starting to update collectorset: %s", newcollectorset.Name)
	_, err := CreateOrUpdateCollectorSet(newcollectorset, c.LogicmonitorClient, c.Clientset)
	if err != nil {
		log.Errorf("Failed to update collectorset: %v", err)
		return
	}

	collectorsetCopy := newcollectorset.DeepCopy()

	if err = c.save(collectorsetCopy); err != nil {
		log.Errorf("Failed to update policy: %v", err)
	}

	log.Infof("Finished updating CollectorSet: %s", collectorsetCopy.Name)
}

func (c *Controller) deleteFunc(obj interface{}) {
	collectorset := obj.(*crv1alpha1.CollectorSet)

	log.Infof("Starting to delete collectorset: %s", collectorset.Name)
	if err := DeleteCollectorSet(collectorset, c.Clientset); err != nil {
		log.Errorf("Failed to delete collectorset: %v", err)
		return
	}

	if err := c.remove(collectorset); err != nil {
		log.Errorf("Failed to remove policy: %v", err)
	}

	log.Infof("Finished deleting CollectorSet: %s", collectorset.Name)
}

// func (c *Controller) listCollectorSets() (*crv1alpha1.CollectorSetList, error) {
// 	collectorsetList := &crv1alpha1.CollectorSetList{}
// 	err := c.RESTClient.Get().
// 		Resource(crv1alpha1.CollectorSetResourcePlural).
// 		Do().
// 		Into(collectorsetList)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to get CollectorSet list: %v", err)
// 	}

// 	return collectorsetList, nil
// }

func (c *Controller) updateCollectorSetStatus(collectorset *crv1alpha1.CollectorSet, ids []int32) (*crv1alpha1.CollectorSet, error) {
	collectorsetCopy := collectorset.DeepCopy()
	collectorsetCopy.Status = crv1alpha1.CollectorSetStatus{
		State: crv1alpha1.CollectorSetStateRegistered,
		IDs:   ids,
	}

	err := c.RESTClient.Put().
		Name(collectorset.ObjectMeta.Name).
		Namespace(collectorset.ObjectMeta.Namespace).
		Resource(crv1alpha1.CollectorSetResourcePlural).
		Body(collectorsetCopy).
		Do().
		Error()

	if err != nil {
		return nil, fmt.Errorf("Failed to update status: %v", err)
	}

	return collectorsetCopy, nil
}

func checkCollectorRegistrationStatus(lmClient *lm.DefaultApi, ids []int32) (bool, error) {
	total := len(ids)
	ready := 0
	for _, id := range ids {
		restResponse, apiResponse, err := lmClient.GetCollectorById(id, "")
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			return false, fmt.Errorf("Failed to get collector: %v", _err)
		}
		if restResponse.Data.Status != 0 {
			ready++
		}
	}
	if total == ready {
		return true, nil
	}

	return false, nil
}

func waitForCollectorsToRegister(lmClient *lm.DefaultApi, ids []int32) error {
	// A collector generates a UUID upon startup, and if the collector
	// container dies and comes back up on a new node, the UUID will be
	// generated again. This causes the backend to think that there are
	// multiple collectors with the same ID, and will not the newly start
	// collector to register itself. Only after 6 minutes of the old collector
	// being down, will the backend allow the newly created collector to
	// register.
	ticker := time.NewTicker(30 * time.Second)
	for c := ticker.C; ; <-c {
		registered, err := checkCollectorRegistrationStatus(lmClient, ids)
		if err != nil {
			return fmt.Errorf("Failed to check collector registration status: %v", err)
		}
		if registered {
			return nil
		}
	}
}

func (c *Controller) save(collectorset *crv1alpha1.CollectorSet) error {
	p := &policy.Policy{}

	switch *collectorset.Spec.Policy.DistibutionStrategy {
	case distributor.RoundRobin:
		p.DistributionStrategy = &roundrobin.Distributor{}
	default:
		return fmt.Errorf("Invalid distribution strategy %q", collectorset.Spec.Policy.DistibutionStrategy.String())
	}

	if err := p.DistributionStrategy.SetIDs(collectorset.Status.IDs); err != nil {
		return err
	}

	if err := c.Storage.SetPolicy(collectorset.Name, p); err != nil {
		return err
	}

	log.Infof("Using distribution strategy %q for collectorset %q", collectorset.Spec.Policy.DistibutionStrategy.String(), collectorset.Name)
	return nil
}

func (c *Controller) remove(collectorset *crv1alpha1.CollectorSet) error {
	return c.Storage.DeletePolicy(collectorset.Name)
}
