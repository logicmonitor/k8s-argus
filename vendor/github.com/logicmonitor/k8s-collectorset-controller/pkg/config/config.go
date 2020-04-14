package config

import (
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	"gopkg.in/yaml.v2"
)

// Config represents the application's configuration file.
type Config struct {
	*Secrets
	ProxyURL string `yaml:"proxy_url"`
}

// Secrets represents the application's sensitive configuration file.
type Secrets struct {
	Account   string `envconfig:"ACCOUNT"`
	ID        string `envconfig:"ACCESS_ID"`
	Key       string `envconfig:"ACCESS_KEY"`
	ProxyUser string `envconfig:"PROXY_USER"`
	ProxyPass string `envconfig:"PROXY_PASS"`
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

	err = envconfig.Process("collectorset-controller", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
