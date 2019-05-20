package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents the application's configuration file.
type Config struct {
	*Secrets
}

// Secrets represents the application's sensitive configuration file.
type Secrets struct {
	Account string `envconfig:"ACCOUNT"`
	ID      string `envconfig:"ACCESS_ID"`
	Key     string `envconfig:"ACCESS_KEY"`
}

// New returns the application configuration specified by the config file.
func New() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("collectorset-controller", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
