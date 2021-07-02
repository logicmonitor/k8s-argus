package config

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config represents the application's configuration file.
// nolint: maligned
type Config struct {
	*Secrets
	ResourceGroupProperties       ResourceGroupProperties `yaml:"resource_group_props"`
	Intervals                     *Intervals              `yaml:"app_intervals"`
	Address                       string                  `yaml:"address"`
	ClusterName                   string                  `yaml:"cluster_name"`
	LogLevel                      *logrus.Level           `yaml:"log_level"`
	DeleteResources               bool                    `yaml:"delete_resources"`
	DisableAlerting               bool                    `yaml:"disable_alerting"`
	ClusterGroupID                int32                   `yaml:"cluster_group_id" refresh_scope:"no"`
	ResourceContainerGroupID      *int32                  `yaml:"resource_group_id,omitempty"`
	ProxyURL                      string                  `yaml:"proxy_url"`
	IgnoreSSL                     bool                    `yaml:"ignore_ssl"`
	OpenMetricsConfig             *OpenmetricsConfig      `yaml:"openmetrics"`
	NumberOfWorkers               *int                    `yaml:"number_of_workers" refresh_scope:"no"`
	EnableNewResourceTree         bool                    `yaml:"enable_new_resource_tree"`
	EnableNamespacesDeletedGroups bool                    `yaml:"enable_namespaces_deleted_groups"`
	RegisterGenericFilter         bool                    `yaml:"register_generic_filter"`
	DeleteArgusPodAfter           *string                 `yaml:"delete_argus_pod_after"`
	DisableResourceMonitoring     []enums.ResourceType    `yaml:"disable_resource_monitoring"`
	TelemetryCronString           *string                 `yaml:"telemetry_cron_string"`
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

// PropOpts made public coz. yaml unmarshaler need it public
type PropOpts struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Override *bool  `yaml:"override,omitempty"`
}

// ResourceGroupProperties represents the properties applied on resource groups
type ResourceGroupProperties struct {
	Cluster     []PropOpts `yaml:"cluster"`
	Pods        []PropOpts `yaml:"pods"`
	Services    []PropOpts `yaml:"services"`
	Deployments []PropOpts `yaml:"deployments"`
	Nodes       []PropOpts `yaml:"nodes"`
	ETCD        []PropOpts `yaml:"etcd"`
	HPA         []PropOpts `yaml:"hpas"`
}

// Intervals represents default and min values for periodic sync, periodic delete and resource cache sycn intervals
type Intervals struct {
	PeriodicSyncInterval      *time.Duration `yaml:"periodic_sync_interval"`
	PeriodicDeleteInterval    *time.Duration `yaml:"periodic_delete_interval"`
	CacheSyncInterval         *time.Duration `yaml:"cache_sync_interval"`
	PeriodicSyncMinInterval   *time.Duration `yaml:"periodic_sync_min_interval" refresh_scope:"no"`
	PeriodicDeleteMinInterval *time.Duration `yaml:"periodic_delete_min_interval" refresh_scope:"no"`
	CacheSyncMinInterval      *time.Duration `yaml:"cache_sync_min_interval" refresh_scope:"no"`
}

// OpenmetricsConfig represents openmetrics configs
type OpenmetricsConfig struct {
	Port *uint16 `yaml:"port" refresh_scope:"no"`
}

type config struct {
	*Config
	mu       sync.Mutex
	hooks    []ConfHook
	hooksrwm sync.RWMutex
}

var conf = &config{
	Config:   nil,
	mu:       sync.Mutex{},
	hooks:    make([]ConfHook, 0),
	hooksrwm: sync.RWMutex{},
}

func InitConfig() error {
	// Application configuration
	if err := Load(); err != nil {
		return err
	}

	w.AddConfigMapHook(Hook{
		Hook: func(key string, value string) {
			if err := conf.UpdateConfig(value); err != nil {
				logrus.Errorf("Failed to hot load config with error: %s", err)
			}
		},
		Predicate: func(action Action, key string, value string) bool {
			return action == Set && key == constants.ConfigFileName
		},
	})
	return nil
}

func (c *config) getConf() *Config {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Config
}

func (c *config) setConf(conf *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prev := c.Config
	c.Config = conf
	go func() {
		c.hooksrwm.RLock()
		defer c.hooksrwm.RUnlock()
		for _, hook := range c.hooks {
			if hook.Predicate(prev, conf) {
				hook.Hook(prev, conf)
			}
		}
	}()
}

// AddConfigHook returns the application configuration specified by the config file.
func (c *config) AddConfigHook(hook ConfHook) {
	c.hooksrwm.Lock()
	defer c.hooksrwm.Unlock()
	c.hooks = append(c.hooks, hook)
	if conf := c.getConf(); hook.Predicate(conf, conf) {
		hook.Hook(conf, conf)
	}
}

// UpdateConfig returns the application configuration specified by the config file.
func (c *config) UpdateConfig(value string) error {
	uconf := &Config{} // nolint: exhaustivestruct

	err := yaml.Unmarshal([]byte(value), uconf)
	if err != nil {
		return err
	}

	err = envconfig.Process(constants.EnvVarArgusConfigPrefix, uconf)
	if err != nil {
		return err
	}

	validateConfig(uconf)
	pconf, err := GetConfig()
	if err != nil {
		return err
	}
	if pconf == nil {
		pconf = uconf
	}

	postLoad(pconf, uconf)

	// update config
	c.setConf(uconf)
	return nil
}

func postLoad(pconf *Config, uconf *Config) {
	logrus.Tracef("Before Old: %v", spew.Sdump(pconf))
	logrus.Tracef("Before New: %v", spew.Sdump(uconf))
	t := reflect.TypeOf(pconf).Elem()
	v := reflect.ValueOf(pconf).Elem()
	nv := reflect.ValueOf(uconf).Elem()
	retainFields(t, v, nv)
	logrus.Tracef("After Old: %v", spew.Sdump(pconf))
	logrus.Tracef("After New: %v", spew.Sdump(uconf))
	// retain refresh scope no fields
}

func retainFields(t reflect.Type, v reflect.Value, nv reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if v.Field(i).Kind() == reflect.Ptr && v.Field(i).Elem().Kind() == reflect.Struct {
			retainFields(f.Type.Elem(), v.Field(i).Elem(), nv.Field(i).Elem())
			continue
		}
		val, ok := f.Tag.Lookup("refresh_scope")
		if ok && val == "no" {
			logrus.Tracef("Field: %s has refresh_scope: %s", f.Name, val)
			if (v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil()) || v.Field(i).Kind() != reflect.Ptr {
				nv.Field(i).Set(v.Field(i))
			}
		}
	}
}

// GetConfig returns the application configuration specified by the config file.
func GetConfig() (*Config, error) {
	conf.mu.Lock()
	defer conf.mu.Unlock()
	return conf.Config, nil
}

func validateConfig(conf *Config) {
	if conf.LogLevel == nil {
		fmt.Print("logLevel is not configured, setting it to default \"info\"") // nolint: forbidigo
		*conf.LogLevel = logrus.InfoLevel
	}
	if conf.ResourceContainerGroupID == nil {
		rootGroup := constants.RootResourceGroupID
		conf.ResourceContainerGroupID = &rootGroup
	}
	if conf.NumberOfWorkers == nil {
		defaultWorkers := 10
		conf.NumberOfWorkers = &defaultWorkers
	}
	if conf.DeleteArgusPodAfter == nil {
		scheduledDeleteTime := "P10DT0H0M0S"
		conf.DeleteArgusPodAfter = &scheduledDeleteTime
	}
	if conf.TelemetryCronString == nil {
		// Defaults to 10 minute if not specified
		defaultTelemetryCron := "*/10 * * * *"
		conf.TelemetryCronString = &defaultTelemetryCron
	}

	validateIntervals(conf.Intervals)
	validateOpenMetricsConfig(conf.OpenMetricsConfig)
}

func validateIntervals(i *Intervals) {
	if i.CacheSyncMinInterval == nil {
		*i.CacheSyncMinInterval = constants.DefaultCacheResyncInterval
	}
	if i.CacheSyncInterval == nil {
		*i.CacheSyncInterval = *i.CacheSyncMinInterval
	} else if *i.CacheSyncInterval < *i.CacheSyncMinInterval {
		*i.CacheSyncInterval = *i.CacheSyncMinInterval
		logrus.Warnf("Please provide valid value for cacheSyncMinInterval. Continuing with default: %v", *i.CacheSyncMinInterval)
	}

	if i.PeriodicSyncMinInterval == nil {
		*i.PeriodicSyncMinInterval = constants.DefaultPeriodicSyncInterval
	}
	if i.PeriodicSyncInterval == nil {
		*i.PeriodicSyncInterval = *i.PeriodicSyncMinInterval
	} else if *i.PeriodicSyncInterval < *i.PeriodicSyncMinInterval {
		*i.PeriodicSyncInterval = *i.PeriodicSyncMinInterval
		logrus.Warnf("Please provide valid value for periodicSyncInterval. Continuing with default: %v", *i.PeriodicSyncInterval)
	}

	if i.PeriodicDeleteMinInterval == nil {
		*i.PeriodicDeleteMinInterval = constants.DefaultPeriodicDeleteInterval
	}
	if i.PeriodicDeleteInterval == nil {
		*i.PeriodicDeleteInterval = *i.PeriodicDeleteMinInterval
	} else if *i.PeriodicDeleteInterval < *i.PeriodicDeleteMinInterval {
		*i.PeriodicDeleteInterval = *i.PeriodicDeleteMinInterval
		logrus.Warnf("Please provide valid value for periodicDeleteInterval. Continuing with default: %v", *i.PeriodicDeleteInterval)
	}
}

func validateOpenMetricsConfig(oc *OpenmetricsConfig) {
	if oc.Port == nil {
		logrus.Warnf("Looks like helm chart is of previous version than the current Argus expects. Please upgrade helm chart. Setting \"%s\" to its default : %v", "openmetrics.port", "2112")
		*oc.Port = 2112
	}
}
