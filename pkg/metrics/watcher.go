package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// InitBulkDiscoveryTime first bulk discovery
var InitBulkDiscoveryTime = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: ArgusNamespace,
		Subsystem: WatcherSubsystem,
		Name:      "init_discovery_time",
		Help:      "First bulk discovery time on start",
	},
	[]string{"resource"},
)

// DeleteEventLatencySummary time taken to receive delete event to argus watcher
var DeleteEventLatencySummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:   ArgusNamespace,
		Subsystem:   WatcherSubsystem,
		Name:        "event_latency",
		Help:        "Time taken to receive event to watcher",
		Objectives:  map[float64]float64{0.25: 0.025, 0.5: 0.05, 0.75: 0.075, 0.9: 0.01, 0.99: 0.001}, // nolint: gomnd
		ConstLabels: map[string]string{"event": "delete"},
	},
	[]string{"resource"},
)

// DeleteEventMissingTimestamp time taken to receive delete event to argus watcher
var DeleteEventMissingTimestamp = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace:   ArgusNamespace,
		Subsystem:   WatcherSubsystem,
		Name:        "missing_timestamp",
		Help:        "Time taken to receive event to watcher",
		ConstLabels: map[string]string{"event": "delete"},
	},
	[]string{"resource"},
)

// ResourceHandlerProcessingTimeSummary event processing time
var ResourceHandlerProcessingTimeSummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  ArgusNamespace,
		Subsystem:  WatcherSubsystem,
		Name:       "processing_time",
		Help:       "Resource watcher event handler processing time for events",
		Objectives: map[float64]float64{0.25: 0.025, 0.5: 0.05, 0.75: 0.075, 0.9: 0.01, 0.99: 0.001}, // nolint: gomnd
	},
	[]string{"resource", "event"},
)
