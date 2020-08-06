package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	argus "github.com/logicmonitor/k8s-argus/pkg"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/healthz"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/permission"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		conf, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Failed to open %s: %v", constants.ConfigPath, err)
			os.Exit(1)
		}

		// Set the logging level.
		if conf.Debug {
			log.SetLevel(log.DebugLevel)
		}
		// starting thread to reflect log levels dynamically
		go func(initLevel bool) {
			t := time.NewTicker(5 * time.Second)
			c := initLevel
			for {
				<-t.C
				conf, err := config.GetConfig()
				if err == nil && c != conf.Debug {
					c = conf.Debug
					if conf.Debug {
						log.Info("Setting debug")
						log.SetLevel(log.DebugLevel)
					} else {
						log.Info("Setting info")
						log.SetLevel(log.InfoLevel)
					}
				}
			}
		}(conf.Debug)

		// Add hook to log pod id in log context
		hook := &lmlog.DefaultFieldHook{}
		log.AddHook(hook)

		// Instantiate the base struct.
		base, err := argus.NewBase(conf)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Init the permission component
		permission.Init(base.K8sClient)

		// Set up a gRPC connection and CSC Client.
		connection.Initialize(conf)

		connection.CreateConnectionHandler()

		// Instantiate the application and add watchers.
		argus, err := argus.NewArgus(base)
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

func init() {
	RootCmd.AddCommand(watchCmd)
}
