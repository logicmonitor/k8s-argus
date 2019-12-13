package utilities

import (
	"regexp"
)

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
