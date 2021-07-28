package eventprocessor

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
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
	go func(inch chan *RunnerCommand) {
		defer func() {
			r.running = false
		}()
		for {
			timeout := time.NewTicker(idleNotifyTimeout)
			select {
			case command := <-inch:
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
	log.Debugf("Executing runner command")
	command.Execute()
	log.Debugf("runner command execution completed")
}
