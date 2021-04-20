package config

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the application's configuration file.
// nolint: maligned
type Config struct {
	*Secrets
	DeviceGroupProperties             DeviceGroupProperties `yaml:"device_group_props"`
	Intervals                         Intervals             `yaml:"app_intervals"`
	Address                           string                `yaml:"address"`
	ClusterCategory                   string                `yaml:"cluster_category"`
	ClusterName                       string                `yaml:"cluster_name"`
	Debug                             bool                  `yaml:"debug"`
	DeleteDevices                     bool                  `yaml:"delete_devices"`
	DisableAlerting                   bool                  `yaml:"disable_alerting"`
	FullDisplayNameIncludeNamespace   bool                  `yaml:"displayName_include_namespace"`
	FullDisplayNameIncludeClusterName bool                  `yaml:"displayName_include_clustername"`
	ClusterGroupID                    int32                 `yaml:"cluster_group_id"`
	ProxyURL                          string                `yaml:"proxy_url"`
	IgnoreSSL                         bool                  `yaml:"ignore_ssl"`
	OpenmetricsConfig                 OpenmetricsConfig     `yaml:"openmetrics"`
}

// Secrets represents the application's sensitive configuration file.
type Secrets struct {
	Account            string `envconfig:"ACCOUNT"`
	ID                 string `envconfig:"ACCESS_ID"`
	Key                string `envconfig:"ACCESS_KEY"`
	EtcdDiscoveryToken string `envconfig:"ETCD_DISCOVERY_TOKEN"`
	ProxyUser          string `envconfig:"PROXY_USER"`
	ProxyPass          string `envconfig:"PROXY_PASS"`
}

// DeviceGroupProperties represents the properties applied on device groups
type DeviceGroupProperties struct {
	Cluster     []map[string]interface{} `yaml:"cluster"`
	Pods        []map[string]interface{} `yaml:"pods"`
	Services    []map[string]interface{} `yaml:"services"`
	Deployments []map[string]interface{} `yaml:"deployments"`
	Nodes       []map[string]interface{} `yaml:"nodes"`
	ETCD        []map[string]interface{} `yaml:"etcd"`
	HPA         []map[string]interface{} `yaml:"hpas"`
}

// Intervals represents default and min values for periodic sync, periodic delete and device cache sycn intervals
type Intervals struct {
	PeriodicSyncInterval      time.Duration `yaml:"periodic_sync_interval"`
	PeriodicDeleteInterval    time.Duration `yaml:"periodic_delete_interval"`
	CacheSyncInterval         time.Duration `yaml:"cache_sync_interval"`
	PeriodicSyncMinInterval   time.Duration `yaml:"periodic_sync_min_interval"`
	PeriodicDeleteMinInterval time.Duration `yaml:"periodic_delete_min_interval"`
	CacheSyncMinInterval      time.Duration `yaml:"cache_sync_min_interval"`
}

// OpenmetricsConfig represents openmetrics configs
type OpenmetricsConfig struct {
	Port int64 `yaml:"port"`
}

// GetConfig returns the application configuration specified by the config file.
func GetConfig() (*Config, error) {
	configBytes, err := ioutil.ReadFile(constants.ConfigPath)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(configBytes, c)
	if err != nil {
		return nil, err
	}

	err = envconfig.Process("argus", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// GetCacheSyncInterval gets cache resync interval
func (config Config) GetCacheSyncInterval() time.Duration {
	cacheInterval := validateAndGetIntervalValue("cache_sync_interval", config.Intervals.CacheSyncInterval, config.Intervals.CacheSyncMinInterval, constants.DefaultCacheResyncInterval)
	log.Debugf("cache_sync_interval - %v ", cacheInterval)
	return cacheInterval
}

// GetPeriodicSyncInterval gets periodic sync interval
func (config Config) GetPeriodicSyncInterval() time.Duration {
	periodicSyncInterval := validateAndGetIntervalValue("periodic_sync_interval", config.Intervals.PeriodicSyncInterval, config.Intervals.PeriodicSyncMinInterval, constants.DefaultPeriodicSyncInterval)
	log.Debugf("periodic_sync_interval - %v ", periodicSyncInterval)
	return periodicSyncInterval
}

// GetPeriodicDeleteInterval gets periodic delete interval
func (config Config) GetPeriodicDeleteInterval() time.Duration {
	periodicDeleteInterval := validateAndGetIntervalValue("periodic_delete_interval", config.Intervals.PeriodicDeleteInterval, config.Intervals.PeriodicDeleteMinInterval, constants.DefaultPeriodicDeleteInterval)
	log.Debugf("periodic_delete_interval - %v ", periodicDeleteInterval)
	return periodicDeleteInterval
}

// ValidateAndGetIntervalValue parses given interval into duration format. Returns default value if any errors.
func validateAndGetIntervalValue(intervalName string, syncInterval, minInterval, defaultValue time.Duration) time.Duration {
	if syncInterval < minInterval {
		log.Warnf("Please provide valid value for %s. Since invalid value is configured, forcefully setting it to default %v. ", intervalName, defaultValue)
		syncInterval = defaultValue
	}

	if syncInterval == 0 || minInterval == 0 {
		log.Warnf("Looks like helm chart is of previous version than the current Argus expects. Please upgrade helm chart. Setting %s to its default : %v", intervalName, defaultValue)
		syncInterval = defaultValue
	}

	return syncInterval
}

// GetOpenmetricsPort gets openmetrics port
func (config Config) GetOpenmetricsPort() string {
	port := validateOpenmetricsConfig(config.OpenmetricsConfig)
	return strconv.FormatInt(port, 10)
}

func validateOpenmetricsConfig(openmetricsConfig OpenmetricsConfig) int64 {
	if openmetricsConfig.Port == 0 {
		log.Warnf("Looks like helm chart is of previous version than the current Argus expects. Please upgrade helm chart. Setting %s to its default : %v", "openmetrics.port", "2112")
		openmetricsConfig.Port = 2112
	}
	return openmetricsConfig.Port
}
