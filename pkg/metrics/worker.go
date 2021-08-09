package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// WorkerCommandsTotal command received total
var WorkerCommandsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: ArgusNamespace,
		Subsystem: WorkerSubsystem,
		Name:      "commands_received_total",
		Help:      "Number of requests received to worker",
	},
	[]string{"worker_id"},
)
