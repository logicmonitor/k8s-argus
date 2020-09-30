package lmlog

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	log "github.com/sirupsen/logrus"
)

// MonitorConfig keeps eye on log level in config with interval of 5 seconds and changes logger level accordingly.
func MonitorConfig() {
	// starting thread to reflect log levels dynamically
	go func(initLevel bool) {
		t := time.NewTicker(5 * time.Second)
		c := initLevel
		for {
			<-t.C
			conf, err := config.GetConfig()
			if err == nil && c != conf.Debug {
				c = conf.Debug
				if conf.Debug {
					log.Info("Setting debug")
					log.SetLevel(log.DebugLevel)
				} else {
					log.Info("Setting info")
					log.SetLevel(log.InfoLevel)
				}
			}
		}
	}(log.GetLevel() == log.DebugLevel)
}
