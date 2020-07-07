package connection

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	collectorsetconstants "github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcConn  *grpc.ClientConn
	cscClient api.CollectorSetControllerClient
	lock      = &sync.Mutex{}
)

// Initialize - it will initialize gRPC connection & csc client
func Initialize(address string) {
	log.Info("Initializing gRPC connection & CSC Client.")
	createConnection(address)
}

func createConnection(address string) {
	var grpcErr error
	grpcConn, grpcErr = grpc.Dial(address, grpc.WithInsecure())
	if grpcErr != nil {
		log.Fatalf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
	}

	var cscErr error
	cscClient, cscErr = WaitForCollectorSetClient(grpcConn)
	if cscErr != nil {
		log.Fatalf("Error while creating CSC client. Error: %v", cscErr.Error())
	}
}

// RecreateConnection - Creates gRPC conn & CSC client, if grpcConn is in TransientFailure or Shutdown state.
func RecreateConnection() {
	if grpcConn != nil {
		state := grpcConn.GetState()
		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			lock.Lock()
			defer lock.Unlock()
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				config, err := config.GetConfig()
				if err != nil {
					fmt.Printf("Failed to open %s: %v", constants.ConfigPath, err)
					os.Exit(1)
				}
				log.Infof("gRPC connection state is \"%v\". Creating new gRPC connection & CSC Client.", state)
				createConnection(config.Address)
			}
		}
	} else {
		lock.Lock()
		defer lock.Unlock()
		if grpcConn == nil {
			config, err := config.GetConfig()
			if err != nil {
				fmt.Printf("Failed to open %s: %v", constants.ConfigPath, err)
				os.Exit(1)
			}
			log.Info("Creating new gRPC connection & CSC Client.")
			createConnection(config.Address)
		}
	}
}

// CloseGRPCConn - Close gRPC connection
func CloseGRPCConn() {
	err := grpcConn.Close()
	if err != nil {
		log.Fatalf("Error while closing gRPC connection. Error: %v", err)
	}
}

// GetCSCClient - returns CSC client
func GetCSCClient() api.CollectorSetControllerClient {
	return cscClient
}

// WaitForCollectorSetClient - wait for collectorset
func WaitForCollectorSetClient(conn *grpc.ClientConn) (api.CollectorSetControllerClient, error) {
	client := api.NewCollectorSetControllerClient(conn)

	ready, err := pollCollectorSetStatus(conn)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !ready {
		log.Fatalf("The collectorset controller does not have any ready collectors")
	}
	log.Infof("The collectorset controller has available collectors")

	return client, nil
}

func pollCollectorSetStatus(conn *grpc.ClientConn) (bool, error) {
	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for collectors to become available")
		case <-ticker.C:
			log.Debugf("Checking collectors status")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			defer cancel()
			req := &healthpb.HealthCheckRequest{
				Service: collectorsetconstants.HealthServerServiceName,
			}
			hc := healthpb.NewHealthClient(conn)
			healthCheckResponse, err := hc.Check(ctx, req)
			if err != nil {
				log.Errorf("Failed to get health check: %v", err)
			}
			if healthCheckResponse.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return true, nil
			}
			log.Debugf("The collectors are not ready: %d", healthCheckResponse.GetStatus())
		}
	}
}
