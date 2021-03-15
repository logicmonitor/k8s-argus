package config

import (
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the application's configuration file.
// nolint: maligned
type Config struct {
	*Secrets
	DeviceGroupProperties             DeviceGroupProperties `yaml:"device_group_props"`
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
