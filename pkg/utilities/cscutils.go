package utilities

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
)

// GetCollectorID - get collectorID from csc
func GetCollectorID() (int32, error) {
	// when argus is running out cluster on local env, grpc connection with csc of cluster cannot be opened hence
	// returns static id
	if IsLocal() {
		id, err := strconv.ParseInt(os.Getenv("COLLECTOR_ID"), 10, 32)
		if err != nil {
			return 0, fmt.Errorf("could not parse collector id from ENV: COLLECTOR_ID: %w", err)
		}

		return int32(id), nil
	}
	reply, err := connection.GetCSCClient().CollectorID(context.Background(), &api.CollectorIDRequest{})
	if err != nil || reply == nil {
		return 0, fmt.Errorf("failed to get collector ID: %w", err)
	}

	return reply.Id, nil
}
