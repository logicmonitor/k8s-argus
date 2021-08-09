package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// APIResponseTimeSummary metric gauge to collect lm api status code count
var APIResponseTimeSummary = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  ArgusNamespace,
		Subsystem:  LMAPISubsystem,
		Name:       "response_time",
		Help:       "LM Rest API response time",
		Objectives: map[float64]float64{0.25: 0.025, 0.5: 0.05, 0.75: 0.075, 0.9: 0.01, 0.99: 0.001}, // nolint: gomnd
	},
	[]string{"url", "method", "code"},
)
