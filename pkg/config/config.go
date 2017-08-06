package config

import (
	"io/ioutil"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	yaml "gopkg.in/yaml.v2"
)

// Config represents the application's configuration file.
type Config struct {
	*Secrets
	ClusterCategory            string `yaml:"cluster_category"`
	ClusterName                string `yaml:"cluster_name"`
	CollectorDescription       string `yaml:"collector_description"`
	CollectorEscalationChainID int32  `yaml:"collector_escalation_chain_id"`
	PreferredCollector         int32
	Debug                      bool `yaml:"debug"`
	DeleteDevices              bool `yaml:"delete_devices"`
	DisableAlerting            bool `yaml:"disable_alerting"`
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
		log.Fatal(err.Error())
	}

	return c, nil
}
