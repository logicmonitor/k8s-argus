package utilities

import (
	"sort"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateSelectorAppliesTo generates selector expression string from given object.
func GenerateSelectorAppliesTo(selector interface{}) string {
	switch val := selector.(type) { //nolint: gocritic
	case *metav1.LabelSelector:
		return coalesceLabelSelectorToAppliesTo(*val)
	case metav1.LabelSelector:
		return coalesceLabelSelectorToAppliesTo(val)
	case map[string]string:
		return coalesceMapToAppliesTo(val)
	default:
		return constants.FalseAppliesTo
	}
}

func coalesceMapToAppliesTo(selector map[string]string) string {
	if len(selector) == 0 {
		return constants.FalseAppliesTo
	}
	sb := strings.Builder{}
	keys := make([]string, 0, len(selector))
	for k := range selector {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sb.WriteString(constants.LabelCustomPropertyPrefix + key + " == \"" + selector[key] + "\"" + constants.LogicalAND)
	}

	return strings.TrimSuffix(sb.String(), constants.LogicalAND)
}

func coalesceLabelSelectorToAppliesTo(labelSelector metav1.LabelSelector) string {
	str := coalesceMapToAppliesTo(labelSelector.MatchLabels)
	if len(labelSelector.MatchExpressions) == 0 {
		return str
	}
	sb := strings.Builder{}
	if str != constants.FalseAppliesTo {
		sb.WriteString(str + constants.LogicalAND)
	}
	for _, expr := range labelSelector.MatchExpressions {
		propKey := constants.LabelCustomPropertyPrefix + expr.Key
		switch expr.Operator {
		case metav1.LabelSelectorOpExists:
			sb.WriteString("exists(\"" + propKey + "\")")
		case metav1.LabelSelectorOpDoesNotExist:
			sb.WriteString("!exists(\"" + propKey + "\")")
		case metav1.LabelSelectorOpIn:
			regex := getAppliesToRegex(expr.Values)
			sb.WriteString(propKey + " =~ \"" + regex + "\"")
		case metav1.LabelSelectorOpNotIn:
			regex := getAppliesToRegex(expr.Values)
			sb.WriteString(propKey + " !~ \"" + regex + "\"")
		}
		sb.WriteString(constants.LogicalAND)
	}
	return strings.TrimSuffix(sb.String(), constants.LogicalAND)
}

func getAppliesToRegex(values []string) string {
	narr := make([]string, len(values))
	for i, v := range values {
		narr[i] = "(?=^" + v + "$)"
	}
	return strings.Join(narr, "|")
}
