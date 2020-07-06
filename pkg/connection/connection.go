package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	grpcLock1 = &sync.Mutex{}
	grpcLock2 = &sync.Mutex{}
	cscLock1  = &sync.Mutex{}
	cscLock2  = &sync.Mutex{}
)

// CreateGRPCConn - Setup a gRPC connection to the collectorset controller.
func CreateGRPCConn(address string) {
	var grpcErr error
	if grpcConn == nil {
		grpcLock1.Lock()
		defer grpcLock1.Unlock()
		if grpcConn == nil {
			log.Infof("Creating new gRPC connection.")
			grpcConn, grpcErr = grpc.Dial(address, grpc.WithInsecure())
			if grpcErr != nil {
				log.Fatalf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
			}
		}
	}

	/* added to check gRPC connection.
	If it is in TransientFailure or in Shutdown state then create new gRPC connection. */
	state := grpcConn.GetState()
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		grpcLock2.Lock()
		defer grpcLock2.Unlock()
		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			log.Infof("gRPC connection state is \"%v\". Creating new gRPC connection.", state)
			grpcConn, grpcErr = grpc.Dial(address, grpc.WithInsecure())
			if grpcErr != nil {
				log.Fatalf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
			}
		}
	}
}

// GetGRPCConn - returns gRPC connection
func GetGRPCConn() *grpc.ClientConn {
	return grpcConn
}

// CloseGRPCConn - Close gRPC connection
func CloseGRPCConn() {
	err := grpcConn.Close()
	if err != nil {
		log.Fatalf("Error while closing gRPC connection. Error: %v", err)
	}
}

// CreateCSCClient - create CSC client
func CreateCSCClient(grpcConn *grpc.ClientConn) {
	var cscErr error
	if cscClient == nil {
		cscLock1.Lock()
		defer cscLock1.Unlock()
		if cscClient == nil {
			log.Infof("Creating new CSC client.")
			cscClient, cscErr = WaitForCollectorSetClient(grpcConn)
			if cscErr != nil {
				log.Fatalf("Error while creating CSC client. Error: %v", cscErr.Error())
			}
		}
	}

	/* added to check gRPC connection.
	If it is not in ready state then create new csc client. */
	if grpcConn != nil {
		state := grpcConn.GetState()
		if state != connectivity.Ready {
			cscLock2.Lock()
			defer cscLock2.Unlock()
			if state != connectivity.Ready {
				log.Infof("gRPC connection state is \"%v\". Creating new CSC client.", state)
				cscClient, cscErr = WaitForCollectorSetClient(grpcConn)
				if cscErr != nil {
					log.Fatalf("Error while creating CSC client. Error: %v", cscErr.Error())
				}
			}
		}
	}
}

// GetCSCClient - returns CSC client
func GetCSCClient() api.CollectorSetControllerClient {
	return cscClient
}

// WaitForCollectorSetClient - wait for collectorset
func WaitForCollectorSetClient(conn *grpc.ClientConn) (api.CollectorSetControllerClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	state := conn.GetState()
	// Wait for connection to be Ready.
	for ; state != connectivity.Ready && conn.WaitForStateChange(ctx, state); state = conn.GetState() {
		log.Infof("Waiting for gRPC")
	}
	if state != connectivity.Ready {
		log.Fatalf("Failed waiting for gRPC to ready, state is %q", state)
	}

	log.Infof("State of gRPC is %q", state)

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
