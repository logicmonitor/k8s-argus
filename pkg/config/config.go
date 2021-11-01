package config

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
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
	NumberOfParallelRunners       *int                    `yaml:"number_of_parallel_runners" refresh_scope:"no"`
	ParallelRunnerQueueSize       *int                    `yaml:"parallel_runner_queue_size" refresh_scope:"no"`
	EnableNewResourceTree         bool                    `yaml:"enable_new_resource_tree"`
	EnableNamespacesDeletedGroups bool                    `yaml:"enable_namespaces_deleted_groups"`
	RegisterGenericFilter         bool                    `yaml:"register_generic_filter"`
	DeleteInfraPodsAfter          *string                 `yaml:"delete_infra_pods_after"`
	DisableResourceMonitoring     []enums.ResourceType    `yaml:"disable_resource_monitoring"`
	DisableResourceAlerting       []enums.ResourceType    `yaml:"disable_resource_alerting"`
	TelemetryCronString           *string                 `yaml:"telemetry_cron_string"`
	SysIpsWaitTimeout             *time.Duration          `yaml:"sys_ips_wait_timeout"`
	EnableProfiling               *bool                   `yaml:"enable_profiling"`
}

func (conf *Config) ShouldDisableAlerting(rt enums.ResourceType) bool {
	for _, e := range conf.DisableResourceAlerting {
		if e == rt {
			return true
		}
	}
	return false
}

func (conf *Config) IsMonitoringDisabled(rt enums.ResourceType) bool {
	for _, d := range conf.DisableResourceMonitoring {
		if d == rt {
			return true
		}
	}

	return false
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
	Raw           map[string][]PropOpts `yaml:",inline"`
	resourceProps map[enums.ResourceType][]PropOpts
}

func (rgp *ResourceGroupProperties) Get(e enums.ResourceType) []PropOpts {
	if props, ok := rgp.resourceProps[e]; ok {
		return props
	}
	return []PropOpts{}
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
	rwmu     sync.RWMutex
	hooks    []ConfHook
	hooksrwm sync.RWMutex
}

var conf = &config{
	Config:   nil,
	rwmu:     sync.RWMutex{},
	hooks:    make([]ConfHook, 0),
	hooksrwm: sync.RWMutex{},
}

func InitConfig() error {
	// Application configuration
	if err := Load(); err != nil {
		return err
	}
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"cm_hook": "config.yaml"}))
	log := lmlog.Logger(lctx)
	w.AddConfigMapHook(Hook{
		Hook: func(key string, value string) {
			log.Tracef("config update hook called: %s", key)
			if err := conf.UpdateConfig(value); err != nil {
				log.Errorf("Failed to hot load config with error: %s", err)
			}
		},
		Predicate: func(action Action, key string, value string) bool {
			log.Tracef("config update hook predicate called. action: %s, key: %s ", action, key)
			return action == Set && key == constants.ConfigFileName
		},
	})
	return nil
}

func (c *config) getConf() *Config {
	c.rwmu.RLock()
	defer c.rwmu.RUnlock()
	if c.Config == nil {
		return nil
	}
	sc := *c.Config
	return &sc
}

func (c *config) setConf(conf *Config) {
	c.rwmu.Lock()
	defer c.rwmu.Unlock()
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
func AddConfigHook(hook ConfHook) {
	conf.addConfigHook(hook)
}

// AddConfigHook returns the application configuration specified by the config file.
func (c *config) addConfigHook(hook ConfHook) {
	c.hooksrwm.Lock()
	defer c.hooksrwm.Unlock()
	c.hooks = append(c.hooks, hook)
	if conf := c.getConf(); hook.Predicate(nil, conf) {
		hook.Hook(nil, conf)
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
	postProcess(uconf)

	validateConfig(uconf)
	pconf := c.getConf()
	if pconf == nil {
		pconf = uconf
	}

	postLoad(pconf, uconf)

	// update config
	c.setConf(uconf)
	return nil
}

func postProcess(uconf *Config) {
	uconf.ResourceGroupProperties.resourceProps = make(map[enums.ResourceType][]PropOpts)
	for k, v := range uconf.ResourceGroupProperties.Raw {
		rt, err := enums.ParseResourceType(k)
		if err == nil {
			uconf.ResourceGroupProperties.resourceProps[rt] = v
		}
	}
}

func postLoad(pconf *Config, uconf *Config) {
	t := reflect.TypeOf(pconf).Elem()
	v := reflect.ValueOf(pconf).Elem()
	nv := reflect.ValueOf(uconf).Elem()
	retainFields(t, v, nv)
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
func GetConfig(lctx *lmctx.LMContext) (*Config, error) {
	v := conf.getConf()
	if v == nil {
		return nil, fmt.Errorf("config is not available: %v", v)
	}
	log := lmlog.Logger(lctx)
	log.Tracef("Config response: %v", v)
	return v, nil
}

func validateConfig(conf *Config) {
	if conf.LogLevel == nil {
		fmt.Print("logLevel is not configured, setting it to default \"info\"") // nolint: forbidigo
		level := logrus.InfoLevel
		conf.LogLevel = &level
	}
	if conf.ResourceContainerGroupID == nil {
		rootGroup := constants.RootResourceGroupID
		conf.ResourceContainerGroupID = &rootGroup
	}
	if conf.NumberOfWorkers == nil {
		defaultWorkers := 10
		conf.NumberOfWorkers = &defaultWorkers
	}
	if conf.NumberOfParallelRunners == nil {
		defaultRunners := 10
		conf.NumberOfParallelRunners = &defaultRunners
	}
	if conf.ParallelRunnerQueueSize == nil {
		defaultQueueSize := 100
		conf.ParallelRunnerQueueSize = &defaultQueueSize
	}
	if conf.DeleteInfraPodsAfter == nil {
		scheduledDeleteTime := "P10DT0H0M0S"
		conf.DeleteInfraPodsAfter = &scheduledDeleteTime
	}
	if conf.TelemetryCronString == nil {
		// Defaults to 10 minute if not specified
		defaultTelemetryCron := "*/10 * * * *"
		conf.TelemetryCronString = &defaultTelemetryCron
	}
	if conf.SysIpsWaitTimeout == nil {
		// Defaults to 5 minute if not specified
		timeout := 5 * time.Minute // nolint: gomnd
		conf.SysIpsWaitTimeout = &timeout
	}
	if conf.EnableProfiling == nil {
		disable := false
		conf.EnableProfiling = &disable
	}

	validateIntervals(conf.Intervals)
	validateOpenMetricsConfig(conf.OpenMetricsConfig)
}

func validateIntervals(i *Intervals) {
	if i.CacheSyncMinInterval == nil {
		d := constants.DefaultCacheResyncInterval
		i.CacheSyncMinInterval = &d
	}
	if i.CacheSyncInterval == nil {
		i.CacheSyncInterval = i.CacheSyncMinInterval
	} else if *i.CacheSyncInterval < *i.CacheSyncMinInterval {
		*i.CacheSyncInterval = *i.CacheSyncMinInterval
		logrus.Warnf("Please provide valid value for cacheSyncMinInterval. Continuing with default: %v", *i.CacheSyncMinInterval)
	}

	if i.PeriodicSyncMinInterval == nil {
		d := constants.DefaultPeriodicSyncInterval
		i.PeriodicSyncMinInterval = &d
	}
	if i.PeriodicSyncInterval == nil {
		i.PeriodicSyncInterval = i.PeriodicSyncMinInterval
	} else if *i.PeriodicSyncInterval < *i.PeriodicSyncMinInterval {
		*i.PeriodicSyncInterval = *i.PeriodicSyncMinInterval
		logrus.Warnf("Please provide valid value for periodicSyncInterval. Continuing with default: %v", *i.PeriodicSyncInterval)
	}

	if i.PeriodicDeleteMinInterval == nil {
		d := constants.DefaultPeriodicDeleteInterval
		i.PeriodicDeleteMinInterval = &d
	}
	if i.PeriodicDeleteInterval == nil {
		i.PeriodicDeleteInterval = i.PeriodicDeleteMinInterval
	} else if *i.PeriodicDeleteInterval < *i.PeriodicDeleteMinInterval {
		*i.PeriodicDeleteInterval = *i.PeriodicDeleteMinInterval
		logrus.Warnf("Please provide valid value for periodicDeleteInterval. Continuing with default: %v", *i.PeriodicDeleteInterval)
	}
}

func validateOpenMetricsConfig(oc *OpenmetricsConfig) {
	if oc.Port == nil {
		logrus.Warnf("Looks like helm chart is of previous version than the current Argus expects. Please upgrade helm chart. Setting \"%s\" to its default : %v", "openmetrics.port", "2112")
		d := uint16(2112) // nolint: gomnd
		oc.Port = &d
	}
}
