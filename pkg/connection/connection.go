package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
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
func Initialize(lctx *lmctx.LMContext, config *config.Config) error {
	log := lmlog.Logger(lctx)
	log.Info("Initializing gRPC connection & CSC Client.")
	appConfig = config
	return createConnection(lctx)
}

func createConnection(lctx *lmctx.LMContext) error {
	conn, grpcErr := createGRPCConnection(lctx)
	if grpcErr != nil {
		return fmt.Errorf("error while creating gRPC connection. Error: %w", grpcErr)
	}

	setGRPCConn(conn)

	client, cscErr := createCSCClient(lctx)
	if cscErr != nil {
		return fmt.Errorf("error while creating Collectorset-controller client. Error: %w", cscErr)
	}

	setCSCClient(client)
	return nil
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

func createGRPCConnection(lctx *lmctx.LMContext) (*grpc.ClientConn, error) {
	log := lmlog.Logger(lctx)
	timeout := time.After(defaultGRPCDialTimeout)

	var gerr error
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for gRPC connection: %w", gerr)
		default:
			conn, gerr := grpcDial()
			if gerr != nil {
				log.Warnf("Error while creating gRPC connection. Error: %s", gerr)
			} else {
				return conn, nil
			}
			time.Sleep(defaultGRPCDialRetryDuration)
		}
	}
}

func grpcDial() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultGRPCDialRetryDuration)
	defer cancel()
	conn, err := grpc.DialContext(ctx, appConfig.Address, grpc.WithBlock(), grpc.WithInsecure())

	return conn, err
}

func createCSCClient(lctx *lmctx.LMContext) (api.CollectorSetControllerClient, error) {
	log := lmlog.Logger(lctx)
	conn := getGRPCConn()
	client := api.NewCollectorSetControllerClient(conn)

	timeout := time.After(defaultCSCClientTimeout)
	hc := healthpb.NewHealthClient(conn)

	var lastKnownStatus healthpb.HealthCheckResponse_ServingStatus
	for {
		select {
		case <-timeout:

			return client, fmt.Errorf("timeout waiting for collectors to become available")
		default:
			lastKnownStatus = getCSCHealth(hc).GetStatus()
			if lastKnownStatus == healthpb.HealthCheckResponse_SERVING {
				return client, nil
			}
			log.Warnf("The collectorset controller is not yet ready to serve argus requests (typically it waits for all collector installations to complete)): %s", lastKnownStatus)
			time.Sleep(defaultCSCClientRetryDuration)
		}
	}
}

func getCSCHealth(hc healthpb.HealthClient) *healthpb.HealthCheckResponse {
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

// RunGrpcHeartBeater - It will create a go routine for handling gRPC connection creation
func RunGrpcHeartBeater() {
	go func() {
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"grpc": "heartbeater"}))
		for {
			time.Sleep(grpCDialRetry)
			checkGRPCState(lctx)
		}
	}()
}

// checkGRPCState - It will check gRPC state & call createConnection if required
func checkGRPCState(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	state := getGRPCConn().GetState()
	if state == connectivity.Shutdown {
		log.Infof("gRPC is in \"%s\" state. Creating new gRPC connection & CSC client.", state)
		if err := createConnection(lctx); err != nil {
			log.Errorf("Failed to reinitialise collectorset-controller client")
		}
	}
}
