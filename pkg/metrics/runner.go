package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// RunnerEventsReceivedCount runner received requests
var RunnerEventsReceivedCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: ArgusNamespace,
		Subsystem: RunnerSubsystem,
		Name:      "events_received_total",
		Help:      "Number of events received to runner",
	},
	[]string{"runner_id"},
)
