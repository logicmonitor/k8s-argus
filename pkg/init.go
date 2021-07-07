package argus

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/eventprocessor"
	"github.com/logicmonitor/k8s-argus/pkg/facade"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/lmexec"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resource"
	"github.com/logicmonitor/k8s-argus/pkg/resourcecache"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup/dgbuilder"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/worker"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

func StartWorkers(lctx *lmctx.LMContext, facade types.LMFacade) error {
	log := lmlog.Logger(lctx)
	if conf, err := config.GetConfig(); err == nil {
		for i := 0; i < *conf.NumberOfWorkers; i++ {
			wc := worker.NewWorker(types.NewWConfig(i))
			log.Debugf("Worker Config %d: %v", i, wc.GetConfig().GetChannel())
			_, err := facade.RegisterWorker(wc)
			if err != nil {
				log.Errorf("Failed to register worker [%d] with error: %s", i, err)
				continue
			}
			wc.Run()
		}
	}
	if facade.Count() == 0 {
		return fmt.Errorf("no worker is runnning")
	}
	log.Infof("Number of workers running: %d", facade.Count())
	return nil
}

func CreateArgus(lctx *lmctx.LMContext, lmClient *client.LMSdkGo) (*Argus, error) {
	clctx := lmlog.LMContextWithFields(lctx, logrus.Fields{"argus": "create"})
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	lmrequester, err := createLMRequester(clctx, lmClient)
	if err != nil {
		return nil, err
	}

	resourceCache := resourcecache.NewResourceCache(lmrequester, *conf.Intervals.CacheSyncInterval)
	// Graceful rebuild
	if resourcegroup.GetClusterGroupProperty(lctx, constants.ResyncCacheProp, lmrequester) == "true" {
		resourceCache.Rebuild(lctx)
		clusterGroupID := util.GetClusterGroupID(lctx, lmrequester)
		clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
		defer resourcegroup.DeleteResourceGroupPropertyByName(clctx, clusterGroupID, &models.EntityProperty{Name: constants.ResyncCacheProp, Value: "true"}, lmrequester)
	}

	resourceCache.Run()

	dgBuilder := &dgbuilder.Builder{}
	resourceGroupManager := &resourcegroup.Manager{
		Builder:       dgBuilder,
		LMRequester:   lmrequester,
		ResourceCache: resourceCache,
	}

	resourceManager := &resource.Manager{ // nolint: exhaustivestruct
		LMRequester:          lmrequester,
		ResourceCache:        resourceCache,
		ResourceGroupManager: resourceGroupManager,
	}

	return NewArgus(lmrequester, resourceManager, resourceCache)
}

func createLMExecutor(lmClient *client.LMSdkGo) *lmexec.LMExec {
	return &lmexec.LMExec{
		LMClient: lmClient,
	}
}

func createLMRequester(lctx *lmctx.LMContext, lmClient *client.LMSdkGo) (*types.LMRequester, error) {
	facadeObj := facade.NewFacade()
	if err := StartWorkers(lctx, facadeObj); err != nil {
		return nil, err
	}
	lmExec := createLMExecutor(lmClient)
	return &types.LMRequester{
		LMFacade:   facadeObj,
		LMExecutor: lmExec,
	}, nil
}

func createRunnerFacade(lctx *lmctx.LMContext) (eventprocessor.RunnerFacade, error) {
	facadeObj := eventprocessor.NewRFacade()
	if err := startRunners(lctx, facadeObj); err != nil {
		return nil, err
	}
	return facadeObj, nil
}

func startRunners(lctx *lmctx.LMContext, facade *eventprocessor.RFacade) error {
	log := lmlog.Logger(lctx)
	if conf, err := config.GetConfig(); err == nil {
		for i := 0; i < *conf.NumberOfParallelRunners; i++ {
			wc := eventprocessor.NewRunner(eventprocessor.NewRunnerConfig(i, *conf.ParallelRunnerQueueSize))
			log.Debugf("Runner Config %d: %v", i, wc.GetConfig().GetChannel())
			_, err := facade.RegisterRunner(wc)
			if err != nil {
				log.Errorf("Failed to register runner for: %d", i)
				continue
			}
			wc.Run()
		}
	}
	if facade.Count() == 0 {
		return fmt.Errorf("no runner is runnning")
	}
	log.Infof("Number of runners running: %d", facade.Count())
	return nil
}
