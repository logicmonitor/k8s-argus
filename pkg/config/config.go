package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config represents the application's configuration file.
type Config struct {
	*Secrets
	Address         string `yaml:"address"`
	ClusterCategory string `yaml:"cluster_category"`
	ClusterGroupID  int32  `yaml:"cluster_group_id"`
	ClusterName     string `yaml:"cluster_name"`
	Debug           bool   `yaml:"debug"`
	DeleteDevices   bool   `yaml:"delete_devices"`
	DisableAlerting bool   `yaml:"disable_alerting"`
}

// Secrets represents the application's sensitive configuration file.
type Secrets struct {
	Account            string `envconfig:"ACCOUNT"`
	ID                 string `envconfig:"ACCESS_ID"`
	Key                string `envconfig:"ACCESS_KEY"`
	EtcdDiscoveryToken string `envconfig:"ETCD_DISCOVERY_TOKEN"`
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
