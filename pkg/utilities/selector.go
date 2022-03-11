package utilities

import (
	"sort"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	expressionSeparator    = ", "
	expressionOpenBracket  = "("
	expressionCloseBracket = ")"
)

// GenerateSelectorExpression generates selector expression string from given object.
// expression adheres the rules defined at https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
// this expression string should work against kubectl's "-l" option.
// for ex: kubectl get pods -l "<expr>" should return the pods on which expr eval's true
func GenerateSelectorExpression(selector interface{}) string {
	switch val := selector.(type) { //nolint: gocritic
	case *metav1.LabelSelector:
		return coalesceLabelSelector(*val)
	case metav1.LabelSelector:
		return coalesceLabelSelector(val)
	case map[string]string:
		return coalesceMap(val)
	default:
		return constants.LabelNullPlaceholder
	}
}

func coalesceMap(selector map[string]string) string {
	if len(selector) == 0 {
		return constants.LabelNullPlaceholder
	}
	sb := strings.Builder{}
	keys := make([]string, 0, len(selector))
	for k := range selector {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sb.WriteString(key + "=" + selector[key] + expressionSeparator)
	}

	return strings.TrimSuffix(sb.String(), expressionSeparator)
}

func coalesceLabelSelector(labelSelector metav1.LabelSelector) string {
	str := coalesceMap(labelSelector.MatchLabels)
	if len(labelSelector.MatchExpressions) == 0 {
		return str
	}
	sb := strings.Builder{}
	if str != constants.LabelNullPlaceholder {
		sb.WriteString(str + expressionSeparator)
	}
	for _, expr := range labelSelector.MatchExpressions {
		switch expr.Operator {
		case metav1.LabelSelectorOpExists:
			sb.WriteString(expr.Key)
		case metav1.LabelSelectorOpDoesNotExist:
			sb.WriteString("!" + expr.Key)
		case metav1.LabelSelectorOpIn:
			sb.WriteString(expr.Key + " in " + expressionOpenBracket +
				strings.Join(expr.Values, expressionSeparator) +
				expressionCloseBracket)
		case metav1.LabelSelectorOpNotIn:
			sb.WriteString(expr.Key + " notin " + expressionOpenBracket +
				strings.Join(expr.Values, expressionSeparator) +
				expressionCloseBracket)
		}
		sb.WriteString(expressionSeparator)
	}
	return strings.TrimSuffix(sb.String(), expressionSeparator)
}
