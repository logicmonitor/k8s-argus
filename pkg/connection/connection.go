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
	"gopkg.in/robfig/cron.v2"
)

var (
	grpcConn  *grpc.ClientConn
	cscClient api.CollectorSetControllerClient
	connLock  sync.RWMutex
	cronLock  sync.RWMutex
	appConfig *config.Config
	counter   int
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
	}
	setGRPCConn(conn)

	client, cscErr := createCSCClient()
	if cscErr != nil {
		log.Errorf("Error while creating gRPC connection. Error: %v", cscErr.Error())
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
		case <-ticker.C:
			conn, grpcErr := grpc.Dial(appConfig.Address, grpc.WithInsecure())
			if grpcErr != nil {
				log.Errorf("Error while creating gRPC connection. Error: %v", grpcErr.Error())
			} else {
				return conn, nil
			}
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
		case <-ticker.C:
			healthCheckResponse := getCSCHealth(hc)
			if healthCheckResponse.GetStatus() == healthpb.HealthCheckResponse_SERVING {
				return client, nil
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

// CreateConnectionCronJob - It will create CronJob for handling gRPC connection creation
func CreateConnectionCronJob() {
	log.Info("Creating cron job for connection handling")
	c := cron.New()
	_, err := c.AddFunc("@every 0h0m10s", func() {
		checkGRPCState()
	})
	if err != nil {
		log.Errorf("Error while creating cron job for connection handling. Error: %v", err)
	}
	c.Start()
}

// checkGRPCState - It will check gRPC state & call createConnection if required
func checkGRPCState() {
	state := getGRPCConn().GetState()
	if state == connectivity.Shutdown || getCounter() > 5 {
		cronLock.Lock()
		defer cronLock.Unlock()
		state := getGRPCConn().GetState()
		if state == connectivity.Shutdown || getCounter() > 5 {
			log.Infof("gRPC is in \"%v\" state. Creating new gRPC connection & CSC client.", state.String())
			createConnection()
			resetCounter()
		}
	} else if state == connectivity.Ready {
		resetCounter()
	} else if state == connectivity.Idle || state == connectivity.Connecting || state == connectivity.TransientFailure {
		incCounter()
	}
}

func getCounter() int {
	return counter
}

func resetCounter() {
	counter = 0
}

func incCounter() {
	counter++
}
