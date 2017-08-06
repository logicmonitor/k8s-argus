package collector

import (
	"fmt"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	log "github.com/sirupsen/logrus"
)

// Controller is the controller manager to LogicMonitor collectors.
type Controller struct {
	*types.Base
}

// NewController instantiates and returns a collector Controller.
func NewController(base *types.Base) *Controller {
	return &Controller{
		Base: base,
	}
}

// Init initializes a collector Controller.
func (cm *Controller) Init() error {
	id, err := cm.discoverCollectorID()
	if err != nil {
		return err
	}
	cm.Config.PreferredCollector = id
	log.Infof("Using collector %d", cm.Config.PreferredCollector)

	return nil
}

func (cm *Controller) discoverCollectorID() (int32, error) {
	attempts := 0
	for {
		restResponse, apiResponse, err := cm.LMClient.GetCollectorList("id", -1, 0, "description:"+cm.Config.CollectorDescription)
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			log.Errorf("Failed to discover collector ID: %v", _err)
		}
		switch restResponse.Data.Total {
		case 0:
			if attempts == 6 {
				return -1, fmt.Errorf("Timeout waiting for collector ID")
			}
			log.Infof("No collector found, waiting 10 seconds...")
			time.Sleep(10 * time.Second)
		case 1:
			return restResponse.Data.Items[0].Id, nil
		default:
			return -1, fmt.Errorf("Found %d collectors with description %q", restResponse.Data.Total, cm.Config.CollectorDescription)
		}
	}
}
