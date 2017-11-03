package server

import (
	"context"
	"fmt"
	"net"

	"github.com/logicmonitor/k8s-collectorset-controller/api"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server responsible for handling collector ID
// requests.
type Server struct {
	*health.Server
	Storage   storage.Storage
	countChan <-chan int
}

// New instantiates and returns a Server.
func New(storage storage.Storage, count <-chan int) *Server {
	srv := &Server{
		Server:    health.NewServer(),
		Storage:   storage,
		countChan: count,
	}

	srv.SetServingStatus(constants.HealthServerServiceName, healthpb.HealthCheckResponse_NOT_SERVING)

	return srv
}

// Run starts the gRPC server.
func (srv *Server) Run() {
	s := grpc.NewServer()
	go srv.listenForPolicyCountChanges()
	healthpb.RegisterHealthServer(s, srv)
	api.RegisterCollectorSetControllerServer(s, srv)
	reflection.Register(s)

	// Start the gRPC server.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", "50000"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s.Serve(lis) // nolint: errcheck
}

// CollectorID implements api.CollectorSetControllerServer. It returns the
// next collector ID to use.
func (srv Server) CollectorID(ctx context.Context, req *api.CollectorIDRequest) (*api.CollectorIDReply, error) {
	return srv.enforcePolicy(req)
}

func (srv *Server) enforcePolicy(req *api.CollectorIDRequest) (*api.CollectorIDReply, error) {
	for policy := range srv.Storage.IterPolicies() {
		if policy.Validated(req) {
			return policy.DistributionStrategy.ID(req)
		}
	}

	return nil, fmt.Errorf("Failed to find a policy that matches the request")
}

func (srv *Server) listenForPolicyCountChanges() {
	for {
		count, ok := <-srv.countChan
		if !ok {
			continue
		}
		log.Infof("Number of available policies: %d", count)
		if count == 0 {
			srv.SetServingStatus(constants.HealthServerServiceName, healthpb.HealthCheckResponse_NOT_SERVING)
			log.Warnf("Server is %q", healthpb.HealthCheckResponse_NOT_SERVING)
		} else {
			srv.SetServingStatus(constants.HealthServerServiceName, healthpb.HealthCheckResponse_SERVING)
			log.Infof("Server is %q", healthpb.HealthCheckResponse_SERVING)
		}
	}
}
