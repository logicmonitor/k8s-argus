package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	collectorsetconstants "github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	"github.com/sirupsen/logrus"
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

const (
	defaultGRPCDialRetryDuration  = 10 * time.Second
	defaultGRPCDialTimeout        = 10 * time.Minute
	defaultCSCClientRetryDuration = 10 * time.Second
	defaultCSCClientTimeout       = 10 * time.Minute
	healthRequestTimeout          = time.Millisecond * 500
	grpCDialRetry                 = 10 * time.Second
)

// Initialize - it will initialize gRPC connection & csc client
func Initialize(config *config.Config) {
	logrus.Info("Initializing gRPC connection & CSC Client.")
	appConfig = config
	createConnection()
}

func createConnection() {
	conn, grpcErr := createGRPCConnection()
	if grpcErr != nil {
		logrus.Errorf("Error while creating gRPC connection. Error: %v", grpcErr.Error())

		return
	}

	setGRPCConn(conn)

	client, cscErr := createCSCClient()
	if cscErr != nil {
		logrus.Errorf("Error while creating gRPC connection. Error: %v", cscErr.Error())

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
	timeout := time.After(defaultGRPCDialTimeout)
	ticker := time.NewTicker(defaultGRPCDialRetryDuration)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for gRPC connection")
		default:
			conn, err := grpcDial()
			if err != nil {
				logrus.Errorf("Error while creating gRPC connection. Error: %v", err.Error())
			} else {
				return conn, nil
			}
			<-ticker.C
		}
	}
}

func grpcDial() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultGRPCDialRetryDuration)
	defer cancel()
	conn, err := grpc.DialContext(ctx, appConfig.Address, grpc.WithBlock(), grpc.WithInsecure())

	return conn, err
}

func createCSCClient() (api.CollectorSetControllerClient, error) {
	conn := getGRPCConn()
	client := api.NewCollectorSetControllerClient(conn)

	timeout := time.After(defaultCSCClientTimeout)
	ticker := time.NewTicker(defaultCSCClientRetryDuration)
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
			logrus.Debugf("The collectors are not ready: %v", healthCheckResponse.GetStatus().String())
			<-ticker.C
		}
	}
}

func getCSCHealth(hc healthpb.HealthClient) *healthpb.HealthCheckResponse {
	logrus.Debug("Checking collectors status")

	ctx, cancel := context.WithTimeout(context.Background(), healthRequestTimeout)
	defer cancel()

	req := &healthpb.HealthCheckRequest{
		Service: collectorsetconstants.HealthServerServiceName,
	}
	healthCheckResponse, err := hc.Check(ctx, req)
	if err != nil {
		logrus.Errorf("Failed to get health check: %v", err)
	}

	return healthCheckResponse
}

// CreateConnectionHandler - It will create a go routine for handling gRPC connection creation
func CreateConnectionHandler() {
	go func() {
		for {
			time.Sleep(grpCDialRetry)
			checkGRPCState()
		}
	}()
}

// checkGRPCState - It will check gRPC state & call createConnection if required
func checkGRPCState() {
	state := getGRPCConn().GetState()
	if state == connectivity.Shutdown {
		logrus.Infof("gRPC is in \"%v\" state. Creating new gRPC connection & CSC client.", state.String())
		createConnection()
	}
}
