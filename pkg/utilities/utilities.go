package utilities

import (
	"reflect"
	"regexp"
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

// Contains determines whether obj is in the target, the type supported by the target is array, slice, map
func Contains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}
