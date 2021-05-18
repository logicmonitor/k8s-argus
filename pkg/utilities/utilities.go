package utilities

import (
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

//LogDeleteEventLatency logs latency of receiving delete event to argus.
func LogDeleteEventLatency(deletionTimestamp *v1.Time, name string) {
	if deletionTimestamp != nil {
		latency := time.Since(deletionTimestamp.Time)
		logrus.Infof("Delete event latency for %s is %v", name, latency)
	} else {
		logrus.Warnf("resource delete time was not present in event context for %s", name)
	}
}
