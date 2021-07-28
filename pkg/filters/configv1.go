package filters

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"gopkg.in/yaml.v3"
)

type ConfigV1 struct {
	Filters map[enums.ResourceType]string `yaml:"filters"`
}

func (confv1 *ConfigV1) ToV2() (*Config, error) {
	conf := &Config{Filters: make(map[enums.ResourceType][]Rule)}
	for k, v := range confv1.Filters {
		rules := make([]Rule, 0)

		if v == "" {
			continue
		}

		for _, expr := range strings.Split(v, "||") {
			rule := &Rule{}
			err := yaml.Unmarshal([]byte(expr), rule)
			if err != nil {
				return nil, fmt.Errorf("failed to parse string %s into rule: %w", expr, err)
			}
			rules = append(rules, *rule)
		}
		conf.Filters[k] = rules
	}
	return conf, nil
}
