package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logicmonitor/k8s-argus/pkg"
	"github.com/logicmonitor/k8s-argus/pkg/constants"

	"github.com/logicmonitor/k8s-argus/pkg/config"
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

		// Instantiate the application and add watchers.
		argus, err := argus.NewArgus(base)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Invoke the watcher.
		argus.Watch()

		// Stay alive
		// TODO: Expose a monitoring endpoint.
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}
