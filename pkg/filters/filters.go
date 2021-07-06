package filters

import (
	"fmt"
	"os"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/logicmonitor/k8s-argus/pkg/config"
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
func Eval(lctx *lmctx.LMContext, resource enums.ResourceType, evaluationParams map[string]interface{}) (bool, error) {
	log := lmlog.Logger(lctx)
	rules, exists := filterConfig.Filters[resource]

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

var filterConfig *Config

// Init package init block so that filterConfig-filterConfig will be loaded on application start
func Init(lctx *lmctx.LMContext) {
	clctx := lmlog.LMContextWithFields(lctx, logrus.Fields{"filter": "init"})
	log := lmlog.Logger(clctx)
	// skip launching filterConfig file read when invoked via go test.
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test.") {
		return
	}
	filterConfig = readFilterConfig(clctx)
	log.Infof("Rule engine loaded with rules: %v", filterConfig)
}

// Config config
type Config struct {
	Filters map[enums.ResourceType][]Rule `yaml:"filters"`
}

func readFilterConfig(lctx *lmctx.LMContext) *Config {
	log := lmlog.Logger(lctx)
	configString, err := config.GetWatchConfig("filters-config.yaml")
	if err != nil {
		log.Errorf("Failed to read FiltersConfig conf file: filters-config.yaml")
	}

	return parseConfig(lctx, []byte(configString))
}

func parseConfig(lctx *lmctx.LMContext, configBytes []byte) *Config {
	log := lmlog.Logger(lctx)
	conf := &Config{} // nolint: exhaustivestruct
	log.Tracef("conf bytes %s ", configBytes)
	err := yaml.Unmarshal(configBytes, conf)
	if err != nil {
		log.Warnf("Couldn't parse filters-config.yaml file: %s", err)
		confv1 := &ConfigV1{}
		err := yaml.Unmarshal(configBytes, confv1)
		if err != nil {
			log.Errorf("Couldn't parse filters-config.yaml file to config version v1: %s", err)
			return conf
		}
		if c, er := confv1.ToV2(); er == nil {
			log.Infof("Filters loaded with v1 version, recommended to change argus-configuration.yaml file into new format")
			return c
		}
		log.Errorf("Failed to convert v1 config into v2 format, change argus-configuration.yaml file into new format")
		return conf
	}
	log.Tracef("Filter conf read: %v", conf)

	return conf
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
func (rule *Rule) Evaluate(parameters map[string]interface{}) (interface{}, error) {
	return govaluate.EvaluableExpression(*rule).Evaluate(parameters)
}

// String string repr
func (rule *Rule) String() string {
	return govaluate.EvaluableExpression(*rule).String()
}
