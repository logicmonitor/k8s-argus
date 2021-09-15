package csc

import (
	"context"
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
)

// GetCollectorID - get collectorID from csc
func GetCollectorID() (int32, error) {
	client := connection.GetCSCClient()
	if client == nil {
		return 0, fmt.Errorf("client is not initialized: %v", client)
	}
	reply, err := client.CollectorID(context.Background(), &api.CollectorIDRequest{})
	if err != nil || reply == nil {
		return 0, fmt.Errorf("failed to get collector ID: %w", err)
	}

	return reply.Id, nil
}
