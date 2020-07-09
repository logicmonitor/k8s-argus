package utilities

import (
	"context"

	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	log "github.com/sirupsen/logrus"
)

// GetCollectorID - get collectorID from csc
func GetCollectorID() int32 {
	reply, err := connection.GetCSCClient().CollectorID(context.Background(), &api.CollectorIDRequest{})
	if err != nil {
		log.Errorf("Failed to get collector ID: %v", err)

		connection.CheckCSCHealthAndRecreateConnection()

		reply, err := connection.GetCSCClient().CollectorID(context.Background(), &api.CollectorIDRequest{})
		if err != nil {
			log.Errorf("Failed to get collector ID: %v", err)
		}

		return reply.Id
	}

	return reply.Id
}
