package filters

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Knetic/govaluate"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SanitiseEvalInput replaces unsupported characters with '_'.
func SanitiseEvalInput(expression string) string {
	expression = strings.ReplaceAll(expression, ".", "_")
	expression = strings.ReplaceAll(expression, "-", "_")

	return expression
}

// Eval evaluates filtering expression based on specified evaluation parameters
func Eval(lctx *lmctx.LMContext, resource enums.ResourceType, evaluationParams govaluate.MapParameters) (bool, error) {
	log := lmlog.Logger(lctx)
	rules, exists := conf.getConf().Filters[resource]

	if !exists {
		return false, nil
	}
	log.Tracef("Exclude rules : %v", rules)

	var combinedErr []error

	for _, rule := range rules {
		evalResult, err := rule.Evaluate(evaluationParams)
		if err != nil {
			combinedErr = append(combinedErr, fmt.Errorf("%w for expression %s", err, rule.String()))

			continue
		}
		if evalResult != nil && evalResult.(bool) {
			return true, nil
		}
	}

	if len(combinedErr) == 0 {
		return false, nil
	}

	arr := make([]string, 0)
	for _, e := range combinedErr {
		arr = append(arr, e.Error())
	}

	return false, fmt.Errorf("evaluation errors: %s", strings.Join(arr, ", "))
}

//  INTERNAL METHODS

type filterHook func(*Config, *Config)

type filterHookPredicate func(*Config, *Config) bool

type FilterHook struct {
	Hook      filterHook
	Predicate filterHookPredicate
}

type filterConfig struct {
	*Config
	mu       sync.Mutex
	hooks    []FilterHook
	hooksrwm sync.RWMutex
}

var conf = &filterConfig{
	Config:   nil,
	mu:       sync.Mutex{},
	hooks:    make([]FilterHook, 0),
	hooksrwm: sync.RWMutex{},
}

func (c *filterConfig) getConf() *Config {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Config
}

func (c *filterConfig) setConf(conf *Config) {
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

// Init package init block so that filterConfig-filterConfig will be loaded on application start
func Init(lctx *lmctx.LMContext) error {
	clctx := lmlog.LMContextWithFields(lctx, logrus.Fields{"filter": "init"})
	log := lmlog.Logger(clctx)
	c, err := readFilterConfig(clctx)
	if err != nil {
		return fmt.Errorf("failed to initialize rule engine for exclusion filters: %w", err)
	}
	conf.setConf(c)
	config.AddConfigMapHook(config.Hook{
		Hook: func(key string, value string) {
			log.Tracef("config update hook called: %s", key)
			if err := conf.UpdateConfig(clctx, value); err != nil {
				log.Errorf("Failed to reload exclusion filters with error: %s", err)
			}
		},
		Predicate: func(action config.Action, key string, value string) bool {
			log.Tracef("config update hook predicate called. action: %s, key: %s ", action, key)
			return action == config.Set && key == constants.FiltersConfigFileName
		},
	})
	log.Infof("Rule engine loaded with rules: %v", conf.getConf())
	return nil
}

// UpdateConfig returns the application configuration specified by the config file.
func (c *filterConfig) UpdateConfig(clctx *lmctx.LMContext, value string) error {
	uconf, err := parseConfig(clctx, []byte(value))
	if err != nil {
		return err
	}

	// update config
	c.setConf(uconf)
	return nil
}

// Config config
type Config struct {
	Filters map[enums.ResourceType][]Rule `yaml:"filters"`
}

func readFilterConfig(lctx *lmctx.LMContext) (*Config, error) {
	configString, err := config.GetWatchConfig(constants.FiltersConfigFileName)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to read FiltersConfig conf: filters-config.yaml")
	}

	return parseConfig(lctx, []byte(configString))
}

func parseConfig(lctx *lmctx.LMContext, configBytes []byte) (*Config, error) {
	log := lmlog.Logger(lctx)
	conf := &Config{} // nolint: exhaustivestruct
	log.Tracef("conf bytes %s ", configBytes)
	err := yaml.Unmarshal(configBytes, conf)
	if err != nil {
		log.Warnf("Couldn't parse filters-config.yaml file: %s", err)
		confv1 := &ConfigV1{}
		err := yaml.Unmarshal(configBytes, confv1)
		if err != nil {
			return conf, fmt.Errorf("couldn't parse filters-config.yaml file to config version v1: %w", err)
		}
		c, er := confv1.ToV2()
		if er != nil {
			return conf, fmt.Errorf("failed to convert v1 config into v2 format, change argus-configuration.yaml file into new format: %w", er)
		}
		log.Warn("Filters loaded with v1 version, recommended to change argus-configuration.yaml file into new format")
		return c, nil

	}
	log.Tracef("Filter conf read: %v", conf)

	return conf, nil
}

// FilterExpression expr
type FilterExpression string

// MarshalText marshal
func (resourceType FilterExpression) MarshalText() ([]byte, error) {
	return []byte(resourceType), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (resourceType *FilterExpression) UnmarshalText(text []byte) error {
	str := string(text)
	str = SanitiseEvalInput(str)
	if strings.Contains(str, "/") {
		str = strings.ReplaceAll(str, "/", "\\/")
	}
	if _, err := govaluate.NewEvaluableExpression(str); err != nil {
		return err
	}
	*resourceType = FilterExpression(str)

	return nil
}

// Rule rule
type Rule govaluate.EvaluableExpression

// UnmarshalText implements encoding.TextUnmarshaler.
func (rule *Rule) UnmarshalText(text []byte) error {
	str := string(text)
	str = SanitiseEvalInput(str)
	if strings.Contains(str, "/") {
		str = strings.ReplaceAll(str, "/", "\\/")
	}
	expr, err := govaluate.NewEvaluableExpression(str)
	if err != nil {
		return err
	}
	*rule = Rule(*expr)

	return nil
}

// Evaluate eval
func (rule *Rule) Evaluate(parameters govaluate.MapParameters) (interface{}, error) {
	return govaluate.EvaluableExpression(*rule).Evaluate(parameters)
}

// String string repr
func (rule *Rule) String() string {
	return govaluate.EvaluableExpression(*rule).String()
}

func AddFilterHook(hook FilterHook) {
	conf.AddFilterHook(hook)
}

func (c *filterConfig) AddFilterHook(hook FilterHook) {
	c.hooksrwm.Lock()
	defer c.hooksrwm.Unlock()
	c.hooks = append(c.hooks, hook)
	if conf := c.getConf(); hook.Predicate(nil, conf) {
		hook.Hook(nil, conf)
	}
}
