package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// StartGaugeTime track
func StartGaugeTime(gauge prometheus.Gauge) (prometheus.Gauge, time.Time) {
	return gauge, time.Now()
}

// SetGauge duration
func SetGauge(gauge prometheus.Gauge, start time.Time) {
	gauge.Set(float64(time.Since(start).Nanoseconds()))
}

// StartTimeObserver track observer
func StartTimeObserver(observer prometheus.Observer) (prometheus.Observer, time.Time) {
	return observer, time.Now()
}

// ObserveTime duration
func ObserveTime(observer prometheus.Observer, start time.Time) {
	observer.Observe(float64(time.Since(start).Nanoseconds()))
}
