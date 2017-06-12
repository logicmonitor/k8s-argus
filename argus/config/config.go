package config

import "github.com/spf13/viper"

// Config represents the application's configuration file.
type Config struct {
	ClusterCategory    string
	ClusterName        string
	Company            string
	ID                 string
	Key                string
	Debug              bool
	DisableAlerting    bool
	PreferredCollector int32
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
	id := viper.GetString("id")
	key := viper.GetString("key")
	preferredCollector := int32(viper.GetInt("preferred_collector"))

	return &Config{
		ClusterCategory:    clusterCategory,
		ClusterName:        clusterName,
		Company:            company,
		Debug:              debug,
		DisableAlerting:    disableAlerting,
		ID:                 id,
		Key:                key,
		PreferredCollector: preferredCollector,
	}
}
