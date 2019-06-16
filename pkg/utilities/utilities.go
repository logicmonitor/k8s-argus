package utilities

import (
	"fmt"
	"regexp"
	"strings"
)

// BuildSystemCategoriesFromLabels formats a system.categories string.
func BuildSystemCategoriesFromLabels(categories string, labels map[string]string) string {
	for k, v := range labels {
		categories += "," + k + "=" + v
	}
	return categories
}

// GetLabelByPrefix takes a list of labels returns the first label matching the specified prefix
func GetLabelByPrefix(prefix string, labels map[string]string) (string, string) {
	for k, v := range labels {
		if match, err := regexp.MatchString("^"+prefix, k); match {
			if err != nil {
				continue
			}
			return k, v
		}
	}
	return "", ""
}

// GetBatchDisplayNames generate a batch of available display names
func GetBatchDisplayNames(baseDisplayName string, clusterName string, count int, startIndex *int) []string {
	var displayNames []string
	reg, err := regexp.Compile("[ ]+")
	if err != nil {
		return displayNames
	}
	formatClusterName := reg.ReplaceAllString(strings.Trim(clusterName, " "), "_")
	for i := 0; i < count; i++ {
		if *startIndex == 0 {
			displayNames = append(displayNames, fmt.Sprintf("%s-%s", baseDisplayName, formatClusterName))
		} else {
			displayNames = append(displayNames, fmt.Sprintf("%s-%s-%d", baseDisplayName, formatClusterName, *startIndex))
		}
		*startIndex++
	}
	return displayNames
}
