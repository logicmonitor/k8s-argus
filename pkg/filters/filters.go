package filters

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	filter        filters
	expressionMap map[string]string
)

// package init block so that filter-config will be loaded on application start
func init() {
	// skip launching config file read when invoked via go test.
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test.") {
		return
	}
	filter = filters{}
	filter.setConfig(readFilterConfig())
}

type filters struct {
	config config
}

type config struct {
	FilterExp filterExpression `yaml:"filter"`
}

type filterExpression struct {
	POD        string `yaml:"pods"`
	SERVICES   string `yaml:"services"`
	DEPLOYMENT string `yaml:"deployments"`
	NODE       string `yaml:"nodes"`
}

func (config config) get(resource string) filterExpression {
	switch resource {
	case "filter":
		return config.FilterExp
	}
	return filterExpression{}
}

func (expression filterExpression) get(resource string) string {
	switch resource {
	case constants.Pods:
		return expression.POD
	case constants.Deployments:
		return expression.DEPLOYMENT
	case constants.Services:
		return expression.SERVICES
	case constants.Nodes:
		return expression.NODE
	}
	return ""
}

// setConfig sets filter config and prepares expression map.
func (f *filters) setConfig(config *config) {
	f.config = *config
	compileExpressionMap()
}

func readFilterConfig() *config {
	configBytes, err := ioutil.ReadFile("/etc/argus/filters-config.yaml")
	if err != nil {
		log.Errorf("Failed to read filters config file: /etc/argus/filters-config.yaml")
	}
	config := &config{}
	log.Debugf("config bytes %s ", configBytes)
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		log.Errorf("Couldn't parse filters-config file.")
	}
	log.Infof("Filter config read: %v", config)
	return config
}

func compileExpressionMap() {
	expressionMap = make(map[string]string)
	expressionMap[constants.Pods] = getFilterExpressionForResource(constants.Pods)
	expressionMap[constants.Deployments] = getFilterExpressionForResource(constants.Deployments)
	expressionMap[constants.Nodes] = getFilterExpressionForResource(constants.Nodes)
	expressionMap[constants.Services] = getFilterExpressionForResource(constants.Services)
}

func getFilterExpressionForResource(resource string) string {
	return filter.config.get("filter").get(resource)
}

func parseFilterExpressions(expression string) []string {
	return strings.Split(expression, "||")
}

// Eval evaluates filtering expression based on specified evaluation parameters
func Eval(resource string, evaluationParams map[string]interface{}) bool {
	filterExpression, exists := expressionMap[resource]

	if !exists {
		log.Debugf("Filtering not possible for resouce %s as entry is missing in configuration.", resource)
		return false
	}

	if len(filterExpression) == 0 {
		log.Debugf("No filtering specified for resouce %s ", resource)
		return false
	}

	if isFilterAll(filterExpression) {
		return true
	}

	parsedExpression := parseFilterExpressions(filterExpression)
	log.Debugf("parsed expression for resource %s: %q", resource, parsedExpression)

	for _, expr := range parsedExpression {
		if strings.Contains(expr, "/") {
			expr = strings.ReplaceAll(expr, "/", "\\/")
		}

		expression, err := govaluate.NewEvaluableExpression(expr)

		if err != nil {
			log.Errorf("Invalid filter expression for resource %s -> %s", resource, expr)
			return false
		}

		result, err := expression.Evaluate(evaluationParams)
		if err != nil {
			log.Debugf("Error while evaluating expression %s", expr)
			continue
		}

		if result.(bool) {
			return true
		}
	}

	return false

}

func isFilterAll(expression string) bool {
	return expression == "*"
}
