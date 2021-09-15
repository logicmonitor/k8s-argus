package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CacheBuilderSummary Cache rebuild time
var CacheBuilderSummary = promauto.NewSummary(
	prometheus.SummaryOpts{
		Namespace:  ArgusNamespace,
		Subsystem:  CacheSubsystem,
		Name:       "build_time",
		Help:       "Time taken to build the cache",
		Objectives: map[float64]float64{0.25: 0.025, 0.5: 0.05, 0.75: 0.075, 0.9: 0.01, 0.99: 0.001}, // nolint: gomnd
	},
)

// SyncTimeSummary synchronize execution time
var SyncTimeSummary = promauto.NewSummary(
	prometheus.SummaryOpts{
		Namespace:  ArgusNamespace,
		Subsystem:  SyncSubsystem,
		Name:       "exec_time",
		Help:       "Time taken to synchronise the resources",
		Objectives: map[float64]float64{0.25: 0.025, 0.5: 0.05, 0.75: 0.075, 0.9: 0.01, 0.99: 0.001}, // nolint: gomnd
	},
)
