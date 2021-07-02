package lmlog

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/sirupsen/logrus"
)

const defaultMonitorDuration = 5 * time.Second

// MonitorConfig keeps eye on log level in config with interval of 5 seconds and changes logger level accordingly.
func MonitorConfig() {
	// starting thread to reflect log levels dynamically
	go func(initLevel logrus.Level) {
		c := initLevel

		t := time.NewTicker(defaultMonitorDuration)
		for range t.C {
			conf, err := config.GetConfig()
			if err == nil && c != *conf.LogLevel {
				logrus.Infof("Setting log level %s", conf.LogLevel)
				logrus.SetLevel(*conf.LogLevel)
				c = *conf.LogLevel
			}
		}
	}(logrus.GetLevel())
}
