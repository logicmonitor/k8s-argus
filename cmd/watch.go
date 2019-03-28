package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/healthz"
	"github.com/logicmonitor/k8s-collectorset-controller/api"
	collectorsetconstants "github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch Kubernetes events",
	Long:  `Monitors a cluster autonomously by watching events and translating them to LogicMonitor objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: All objects created by the Watcher should add a property that
		// indicates the Watcher version. This can be useful in migrations from one
		// version to the next.

		// Application configuration
		config, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Failed to open %s: %v", constants.ConfigPath, err)
			os.Exit(1)
		}

		// Set the logging level.
		if config.Debug {
			log.SetLevel(log.DebugLevel)
		}

		// Instantiate the base struct.
		base, err := argus.NewBase(config)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Set up a gRPC connection to the collectorset controller.
		conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err.Error())
		}
		defer conn.Close() // nolint: errcheck
		client, err := waitForCollectorSetClient(conn)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Instantiate the application and add watchers.
		argus, err := argus.NewArgus(base, client)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Invoke the watcher.
		argus.Watch()

		// Health check.
		http.HandleFunc("/healthz", healthz.HandleFunc)

		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func waitForCollectorSetClient(conn *grpc.ClientConn) (api.CollectorSetControllerClient, error) {
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

func init() {
	RootCmd.AddCommand(watchCmd)
}
