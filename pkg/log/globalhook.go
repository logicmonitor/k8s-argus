package lmlog

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// DefaultFieldHook global hook to add pod id in log stmt
type DefaultFieldHook struct {
}

// Fire log hook method before dump
func (hook *DefaultFieldHook) Fire(entry *log.Entry) error {
	entry.Data["pod_id"] = os.Getenv("MY_POD_NAME")
	return nil
}

// Levels level array to enable hook on
func (hook *DefaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}
