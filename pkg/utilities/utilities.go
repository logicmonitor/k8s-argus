package utilities

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
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

// GetShortUUID returns short ids. introduced this util function to start for traceability of events and its logs
func GetShortUUID() uint32 {
	return uuid.New().ID()
}

// GetResourceDisplayName get resource name
func GetResourceDisplayName(resourceType, resourceName, namespace, UID string) string {
	switch resourceType {
	case constants.Pods:
		return fmt.Sprintf("%s-%s", resourceName, namespace)
	case constants.Services:
		return fmt.Sprintf("%s.%s.svc-%s", resourceName, namespace, UID)
	case constants.Deployments:
		return fmt.Sprintf("%s.%s.deploy-%s", resourceName, namespace, UID)
	case constants.Nodes:
		return resourceName
	}
	return fmt.Sprintf("%s-%s", resourceName, namespace)
}
