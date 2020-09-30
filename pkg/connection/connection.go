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
	"google.golang.org/grpc/connectivity"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcConn  *grpc.ClientConn
	cscClient api.CollectorSetControllerClient
	connLock  sync.RWMutex
	appConfig *config.Config
)

// Initialize - it will initialize gRPC connection & csc client
func Initialize(config *config.Config) {
	log.Info("Initializing gRPC connection & CSC Client.")
	appConfig = config
	createConnection()
}

func createConnection() {
	conn, grpcErr := createGRPCConnection()
	if grpcErr != nil {
		log.Errorf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
		return
	}
	setGRPCConn(conn)

	client, cscErr := createCSCClient()
	if cscErr != nil {
		log.Errorf("Error while creating gRPC connection. Error: %v", cscErr.Error())
		return
	}
	setCSCClient(client)
}

func setGRPCConn(conn *grpc.ClientConn) {
	connLock.Lock()
	defer connLock.Unlock()
	grpcConn = conn
}

func getGRPCConn() *grpc.ClientConn {
	connLock.RLock()
	defer connLock.RUnlock()
	return grpcConn
}

func setCSCClient(csc api.CollectorSetControllerClient) {
	connLock.Lock()
	defer connLock.Unlock()
	cscClient = csc
}

// GetCSCClient - returns CSC client
func GetCSCClient() api.CollectorSetControllerClient {
	connLock.RLock()
	defer connLock.RUnlock()
	return cscClient
}

func createGRPCConnection() (*grpc.ClientConn, error) {
	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for gRPC connection")
		default:
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			defer cancel()
			conn, err := grpc.DialContext(ctx, appConfig.Address, grpc.WithBlock(), grpc.WithInsecure())
			if err != nil {
				log.Errorf("Error while creating gRPC connection. Error: %v", err.Error())
			} else {
				return conn, nil
			}
			<-ticker.C
		}
	}
}

func createCSCClient() (api.CollectorSetControllerClient, error) {
	conn := getGRPCConn()
	client := api.NewCollectorSetControllerClient(conn)

	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	hc := healthpb.NewHealthClient(conn)
	for {
		select {
		case <-timeout:
			return client, fmt.Errorf("timeout waiting for collectors to become available")
		default:
			healthCheckResponse := getCSCHealth(hc)
			if healthCheckResponse.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return client, nil
			}
			log.Debugf("The collectors are not ready: %v", healthCheckResponse.GetStatus().String())
			<-ticker.C
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

// CreateConnectionHandler - It will create a go routine for handling gRPC connection creation
func CreateConnectionHandler() {
	go func() {
		for {
			time.Sleep(time.Duration(10) * time.Second)
			checkGRPCState()
		}
	}()
}

// checkGRPCState - It will check gRPC state & call createConnection if required
func checkGRPCState() {
	state := getGRPCConn().GetState()
	if state == connectivity.Shutdown {
		log.Infof("gRPC is in \"%v\" state. Creating new gRPC connection & CSC client.", state.String())
		createConnection()
	}
}
