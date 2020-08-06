package worker

import (
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
)

const (
	// RETRYLIMIT Default retry limit for workers
	RETRYLIMIT = 2
)

// Worker object
type Worker struct {
	config      *types.WConfig
	initialized bool
	running     bool
}

// GetConfig returns the config object
func (w *Worker) GetConfig() *types.WConfig {
	return w.config
}

// NewWorker creates a worker with provided config
func NewWorker(c *types.WConfig) *Worker {
	return &Worker{
		config:      c,
		initialized: false,
		running:     false,
	}
}

// NewSharedReaderWorker Creates new shared reader worker to expedites reads using dedicated parallel channel for reads
func NewSharedReaderWorker(SharedReadChannel chan types.ICommand) *Worker {
	inChan := make(chan types.ICommand)
	w := &Worker{
		config: &types.WConfig{
			MethodChannels: map[string]chan types.ICommand{
				"GET":    SharedReadChannel,
				"POST":   inChan,
				"DELETE": inChan,
				"PUT":    inChan,
				"PATCH":  inChan,
			},
		},
	}
	return w
}

// NewHTTPWorker creates new httpworker with single input channel
func NewHTTPWorker() *Worker {
	inChan := make(chan types.ICommand)
	w := &Worker{
		config: &types.WConfig{
			MethodChannels: map[string]chan types.ICommand{
				"GET":    inChan,
				"POST":   inChan,
				"DELETE": inChan,
				"PUT":    inChan,
				"PATCH":  inChan,
			},
		},
	}
	return w
}

// Run creates go routine to handle requests
func (w *Worker) Run() {
	lctx := lmctx.WithLogger(logrus.WithFields(logrus.Fields{"worker": w.config.ID}))
	log := lctx.Logger()
	log.Infof("Starting worker for %v", w.config.ID)
	channelMap := make(map[chan types.ICommand]bool)
	for _, ch := range w.config.MethodChannels {
		channelMap[ch] = true
	}
	for ch := range channelMap {
		go func(inch chan types.ICommand) {
			for {
				timeout := time.Tick(10 * time.Second)
				select {
				case command := <-inch:
					lctx := command.LMContext()
					commandCtx := lctx.WithFields(logrus.Fields{"worker": w.config.ID})
					log = commandCtx.Logger()
					log.Debugf("Received command")
					w.handleCommand(commandCtx, command)
				case <-timeout:
					log.Debugf("%v worker is idle", w.config.ID)
				}
			}
			// nolint: vet
			w.running = false
		}(ch)
	}
	w.initialized = true
	w.running = true
}

func (w *Worker) handleCommand(lctx *lmctx.LMContext, command types.ICommand) {
	log := lctx.Logger()
	var resp interface{}
	var err error
	retryMax := w.config.RetryLimit
	if retryMax < 1 {
		retryMax = RETRYLIMIT
	}
	for i := 1; i <= retryMax; i++ {
		resp, err = command.Execute()
		if err == nil {
			break
		}
		log.Warningf("Request failed with error %v, retrying for %v time...", err, i)
		time.Sleep(5 * time.Second)
	}

	wresp := &types.WorkerResponse{
		Response: resp,
		Error:    err,
	}
	if n, ok := command.(types.Responder); ok {
		log.Debugf("Sync command")
		if rch := n.GetResponseChannel(); rch != nil {
			log.Debugf("Response channel defined")
			rch <- wresp
		}
	}
}
