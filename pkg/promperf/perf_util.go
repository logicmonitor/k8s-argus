package promperf

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Track track
func Track(gauge prometheus.Gauge) (prometheus.Gauge, time.Time) {
	return gauge, time.Now()
}

// Duration duration
func Duration(gauge prometheus.Gauge, start time.Time) {
	gauge.Set(time.Since(start).Seconds())
	logrus.Infof("trace duration: %v: %v", gauge.Desc(), time.Since(start).Seconds())
}
