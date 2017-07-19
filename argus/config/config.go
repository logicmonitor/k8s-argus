package config

import (
	"github.com/spf13/viper"
)

// Config represents the application's configuration file.
type Config struct {
	ClusterCategory      string
	ClusterName          string
	CollectorDescription string
	Company              string
	ID                   string
	Key                  string
	Debug                bool
	DisableAlerting      bool
	DeleteDevices        bool
	PreferredCollector   int32
}

// GetConfig returns the application configuration specified by environment
// variables.
func GetConfig() *Config {
	// Prefix environment variables with "WATCHER_"
	viper.SetEnvPrefix("watcher")

	clusterCategory := viper.GetString("cluster_category")
	clusterName := viper.GetString("cluster_name")
	company := viper.GetString("company")
	debug := viper.GetBool("debug")
	disableAlerting := viper.GetBool("disable_alerting")
	deleteDevices := viper.GetBool("delete_devices")
	id := viper.GetString("id")
	key := viper.GetString("key")
	collectorDescription := viper.GetString("collector_description")

	return &Config{
		ClusterCategory:      clusterCategory,
		ClusterName:          clusterName,
		CollectorDescription: collectorDescription,
		Company:              company,
		Debug:                debug,
		DisableAlerting:      disableAlerting,
		DeleteDevices:        deleteDevices,
		ID:                   id,
		Key:                  key,
		// PreferredCollector: preferredCollector,
	}
}
