package eventprocessor

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

// idleNotifyTimeout worker idle timeout
var idleNotifyTimeout = 10 * time.Second

// Runner object
type Runner struct {
	config      *RunnerConfig
	initialized bool
	running     bool
}

// NewRunner creates a worker with provided config
func NewRunner(c *RunnerConfig) *Runner {
	return &Runner{
		config:      c,
		initialized: false,
		running:     false,
	}
}

func (r *Runner) GetConfig() *RunnerConfig {
	return r.config
}

func (r *Runner) Run() {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"runner": r.config.ID}))
	log := lmlog.Logger(lctx)
	log.Debugf("Starting Runner")
	g := promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   metrics.ArgusNamespace,
			Subsystem:   metrics.RunnerSubsystem,
			Name:        "queued_count",
			Help:        "Number of fetched events present in runner queue",
			ConstLabels: map[string]string{"runner_id": fmt.Sprintf("%d", r.config.ID)},
		},
		func() float64 {
			return float64(len(r.config.inCh))
		},
	)

	go func(inch chan *RunnerCommand) {
		defer func() {
			r.running = false
			prometheus.Unregister(g)
		}()
		timeout := time.NewTicker(idleNotifyTimeout)
		defer timeout.Stop()
		for {
			select {
			case command := <-inch:
				metrics.RunnerEventsReceivedCount.WithLabelValues(fmt.Sprintf("%d", r.config.ID)).Inc()
				lctx2 := command.Lctx
				commandCtx := lmlog.LMContextWithFields(lctx2, logrus.Fields{"runner": r.config.ID})
				log2 := lmlog.Logger(commandCtx)
				log2.Debugf("Received runner command")
				r.handleCommand(commandCtx, command)
			case <-timeout.C:
				// log.Tracef("%v runner is idle", r.config.ID)
			}
		}
	}(r.config.GetChannel())

	r.initialized = true
	r.running = true
}

func (r *Runner) handleCommand(lctx *lmctx.LMContext, command *RunnerCommand) {
	log := lmlog.Logger(lctx)
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Panic for %s: %s", util.GetCurrentFunctionName(), r)
			log.Errorf("Panic stack trace: %s", debug.Stack())
		}
	}()
	log.Debugf("Executing runner command")
	command.Execute()
	log.Debugf("runner command execution completed")
}
