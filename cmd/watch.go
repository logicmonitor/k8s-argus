package cmd

import (
	"fmt"
	"net/http"
	"os"

	argus "github.com/logicmonitor/k8s-argus/pkg"
	"github.com/logicmonitor/k8s-argus/pkg/client/logicmonitor"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/connection"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/cronjob"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/healthz"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{ // nolint: exhaustivestruct
	Use:   "watch",
	Short: "Watch Kubernetes events",
	Long:  `Monitors a cluster autonomously by watching events and translating them to LogicMonitor objects.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: All objects created by the Watcher should add a property that
		//  indicates the Watcher version. This can be useful in migrations from one
		//  version to the next.

		// Add hook to log pod id and goroutine id in log context
		hook := &lmlog.DefaultFieldHook{}
		logrus.AddHook(hook)

		if kubeConfigFile != "" {
			err := os.Setenv(constants.IsLocal, "true")
			if err != nil {
				logrus.Errorf("Failed to set IsLocal environment")
			}
		}

		cronjob.Init()
		fmt.Printf("kubeconfig file path: %s\n", kubeConfigFile) // nolint: forbidigo
		if err := config.Init(kubeConfigFile); err != nil {
			fmt.Println("Failed to initialise Kubernetes client: %w", err) // nolint: forbidigo
			os.Exit(constants.ConfigInitK8sClientExitCode)
		}

		if err := config.InitConfig(); err != nil {
			fmt.Println("failed to load application config from configmaps") // nolint: forbidigo
			os.Exit(constants.ConfigInitExitCode)
		}
		conf, err := config.GetConfig()
		if err != nil {
			fmt.Println("Failed to get config: %w", err) // nolint: forbidigo
			os.Exit(constants.GetConfigExitCode)
		}

		// Once minimal configuration gets loaded, start config watcher to watch on events
		config.Run()

		if logLevel := os.Getenv("LOG_LEVEL"); util.IsLocal() && logLevel != "" {
			level, err := logrus.ParseLevel(logLevel)
			if err != nil {
				fmt.Println("Incorrect log level, setting to info") // nolint: forbidigo
				logrus.SetLevel(logrus.InfoLevel)
			} else {
				logrus.SetLevel(level)
			}
		} else {
			logrus.SetLevel(*conf.LogLevel)
			// Monitor config for log level change
			registerLogLevelChangeHook()
		}

		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"watch": "init"}))
		log := lmlog.Logger(lctx)

		filters.Init(lctx)
		// LogicMonitor API client.
		lmClient, err := logicmonitor.NewLMClient(conf)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		if util.IsLocal() {
			logrus.SetFormatter(&logrus.TextFormatter{ // nolint: exhaustivestruct
				ForceColors: false,
			})
		} else {
			// Set up a gRPC connection and CSC Client.
			if err := connection.Initialize(lctx, conf); err != nil {
				log.Fatalf("failed to initialize collectorset-controller connection: %s", err)
			}
			connection.RunGrpcHeartBeater()
		}

		argusObj, err := argus.CreateArgus(lctx, lmClient)
		if err != nil {
			return
		}
		err = argusObj.Init()
		// Instantiate the application and add watchers.
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if err := argusObj.CreateWatchers(lctx); err != nil {
			log.Fatal(err.Error())
			return
		}
		argusObj.Watch(lctx)

		// To update K8s & Helm properties in cluster resource group periodically with the server
		err = cronjob.StartTelemetryCron(argusObj.ResourceCache, argusObj.LMRequester)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		http.Handle("/metrics", promhttp.Handler())
		go func() {
			addr := fmt.Sprintf(":%d", *conf.OpenMetricsConfig.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()

		// Health check.
		http.HandleFunc("/healthz", healthz.HandleFunc)

		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

// nolint: gochecknoinits
func init() {
	RootCmd.AddCommand(watchCmd)
}

// registerLogLevelChangeHook keeps eye on log level in config with interval of 5 seconds and changes logger level accordingly.
func registerLogLevelChangeHook() {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"config_hook": "log_level"}))
	log := lmlog.Logger(lctx)
	config.AddConfigHook(config.ConfHook{
		Hook: func(prev *config.Config, updated *config.Config) {
			log.Infof("Setting log level %s", *updated.LogLevel)
			logrus.SetLevel(*updated.LogLevel)
		},
		Predicate: func(prev *config.Config, updated *config.Config) bool {
			return prev == nil || *prev.LogLevel != *updated.LogLevel
		},
	})
}
