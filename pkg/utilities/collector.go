package utilities

import (
	"context"
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	log "github.com/sirupsen/logrus"
)

// GetCollectorID - get collectorID from csc
func GetCollectorID() int32 {
	reply, err := connection.GetCSCClient().CollectorID(context.Background(), &api.CollectorIDRequest{})
	if err != nil || reply == nil {
		log.Errorf("Failed to get collector ID: %v", err)

		// Application configuration
		config, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Failed to open %s: %v", constants.ConfigPath, err)
			os.Exit(1)
		}

		// If collectorset-controller pod restarted then recreate gRPC connection & CSC client.
		connection.CreateGRPCConn(config.Address)
		connection.CreateCSCClient(connection.GetGRPCConn())

		/* recursive call to handle cscClient error like
		'Error while dialing dial tcp 10.97.80.18:50000: connect: connection refused'
		'Failed to find a policy that matches the request'. */
		collectorID := GetCollectorID()
		// reply, err = cscClient.CollectorID(context.Background(), &api.CollectorIDRequest{})

		return collectorID
	}

	return reply.Id
}
