package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	collectorsetconstants "github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcConn  *grpc.ClientConn
	cscClient api.CollectorSetControllerClient
	lock      sync.Mutex
	appConfig *config.Config
)

// Initialize - it will initialize gRPC connection & csc client
func Initialize(config *config.Config) {
	log.Info("Initializing gRPC connection & CSC Client.")
	appConfig = config
	createGRPCConnection()
	createCSCClient()
}

func createGRPCConnection() {
	var grpcErr error
	grpcConn, grpcErr = grpc.Dial(appConfig.Address, grpc.WithInsecure())
	if grpcErr != nil {
		log.Fatalf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
	}
}

// GetCSCClient - returns CSC client
func GetCSCClient() api.CollectorSetControllerClient {
	return cscClient
}

func createCSCClient() {
	cscClient = api.NewCollectorSetControllerClient(grpcConn)

	ready, err := pollCollectorSetStatus()
	if err != nil {
		log.Fatal(err.Error())
	}

	if !ready {
		log.Fatalf("The collectorset controller does not have any ready collectors")
	}
	log.Infof("The collectorset controller has available collectors")
}

func pollCollectorSetStatus() (bool, error) {
	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	hc := healthpb.NewHealthClient(grpcConn)
	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for collectors to become available")
		case <-ticker.C:
			healthCheckResponse := getCSCHealth(hc)
			if healthCheckResponse.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return true, nil
			}
			log.Debugf("The collectors are not ready: %v", healthCheckResponse.GetStatus().String())
		}
	}
}

func getCSCHealth(hc healthpb.HealthClient) *healthpb.HealthCheckResponse {
	log.Debug("Checking collectors status")
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	req := &healthpb.HealthCheckRequest{
		Service: collectorsetconstants.HealthServerServiceName,
	}
	healthCheckResponse, err := hc.Check(ctx, req)
	if err != nil {
		log.Errorf("Failed to get health check: %v", err)
	}
	return healthCheckResponse
}

// CheckCSCHealthAndRecreateConnection - check CSC health if it is not SERVING then recreate gRPC connection & CSC client.
func CheckCSCHealthAndRecreateConnection() {
	hc := healthpb.NewHealthClient(grpcConn)
	healthCheckResponse := getCSCHealth(hc)
	if healthCheckResponse.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		lock.Lock()
		defer lock.Unlock()
		if healthCheckResponse.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			log.Infof("CSC client is in \"%v\" state. Creating new gRPC connection & CSC client.", healthCheckResponse.GetStatus().String())
			createGRPCConnection()
			createCSCClient()
		}
	}
}
