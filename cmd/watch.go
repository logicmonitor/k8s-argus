package cmd

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/logicmonitor/k8s-argus/argus"
	"github.com/logicmonitor/k8s-argus/argus/config"
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
		config := config.GetConfig()
		// Set the logging level.
		if config.Debug {
			log.SetLevel(log.DebugLevel)
		}
		// Instantiate the application and add watchers.
		argus, err := argus.NewArgus(config)
		if err != nil {
			log.Fatal(err.Error())
		}
		// Invoke the watcher and store the response.
		argus.Watch()
		// Stay alive
		// TODO: Expose a monitoring endpoint.
		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}
