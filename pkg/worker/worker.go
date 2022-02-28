package worker

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	// retryLimit Default retry limit for workers
	retryLimit = 2
	// maxRateLimitRetry number of retries to give up on rate limit retry request
	maxRateLimitRetry = 2
	// retryBackoffTimeDurationDefault default
	retryBackoffTimeDurationDefault = 5 * time.Second
	// workerIdleNotifyTimeout worker idle timeout
	workerIdleNotifyTimeout = 10 * time.Second
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

// Run creates go routine to handle requests
func (w *Worker) Run() {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"worker": w.config.ID}))
	log := lmlog.Logger(lctx)
	log.Debugf("Starting worker")
	workerQueueCount := promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   metrics.ArgusNamespace,
			Subsystem:   metrics.WorkerSubsystem,
			Name:        "commands_queued",
			Help:        "Number of worker commands queued",
			ConstLabels: map[string]string{"worker_id": fmt.Sprintf("%d", w.config.ID)},
		},
		func() float64 {
			return float64(len(w.config.GetChannel()))
		},
	)
	go func(inch chan *types.WorkerCommand) {
		defer func() {
			w.running = false
			prometheus.Unregister(workerQueueCount)
		}()
		timeout := time.NewTicker(workerIdleNotifyTimeout)
		defer timeout.Stop()
		for {
			select {
			case command := <-inch:
				metrics.WorkerCommandsTotal.WithLabelValues(fmt.Sprintf("%d", w.config.ID)).Inc()
				lctx2 := command.Lctx
				commandCtx := lmlog.LMContextWithFields(lctx2, logrus.Fields{"worker": w.config.ID})
				log2 := lmlog.Logger(commandCtx)
				log2.Debugf("Received command")
				w.handleCommand(commandCtx, command)
			case <-timeout.C:
				// log.Tracef("%v worker is idle", w.config.ID)
			}
		}
	}(w.config.GetChannel())

	w.initialized = true
	w.running = true
}

func (w *Worker) popRLToken(command *types.WorkerCommand) {
	_, err := Allow(command.APIInfo.GetPatternKey(), true)
	if err != nil {
		return
	}
}

func (w *Worker) handleCommand(lctx *lmctx.LMContext, command *types.WorkerCommand) {
	log := lmlog.Logger(lctx)
	log.Tracef("Poping token")
	w.popRLToken(command)
	log.Tracef("Token popped")
	retryMax := w.config.MaxRetry
	if retryMax < 1 {
		retryMax = retryLimit
	}
	log.Tracef("Executing command")
	resp, err := w.executeWithRetry(lctx, retryMax, command)
	log.Tracef("Command response: %s ; %s", resp, err)
	wresp := &types.WorkerResponse{
		Response: resp,
		Error:    err,
	}
	if command.IsSync() {
		rch := command.GetResponseChannel()
		log.Tracef("Sending response")
		select {
		case rch <- wresp:
		// to make non blocking, golang unbuffered channels are blocking if there is no goroutine listening on channel. if facade timed out and returned then this would become blocking
		case <-command.GetContext().Done():
			log.Warnf("Response cannot sent. ApiInfo: %v", command.APIInfo)
		}
	} else {
		log.Infof("Async command response: %v", wresp)
	}
}

func (w *Worker) executeWithRetry(lctx *lmctx.LMContext, retry int, command *types.WorkerCommand) (interface{}, error) {
	log := lmlog.Logger(lctx)
	var resp interface{}
	var err error

	rateLimitRetry := 0
outer:
	for i := 1; i <= retry && rateLimitRetry < maxRateLimitRetry; i++ {
		resp, err = command.Execute()
		if err == nil {
			break
		}
		code := util.GetHTTPStatusCodeFromLMSDKError(err)
		log.Debugf("Status code: %v", code)

		// retry only if request failed because of server error
		switch {
		// 5XX server error
		case code/100 == 5: // nolint: gomnd
			log.Warningf("Request failed with error %v, retrying for %v time...", err, i)
			time.Sleep(retryBackoffTimeDurationDefault)
		case code == http.StatusTooManyRequests:
			log.Infof("Rate limits reached for: %s", command.APIInfo.GetPatternKey())
			req := w.getRLLimit(lctx, err)
			if req != nil {
				w.setNewLimit(lctx, req, command)
				log.Infof("Waiting for rate limit window %v", req.Window)
				time.Sleep(time.Duration(req.Window) * time.Second)
			} else {
				log.Warnf("Rate limit could not parse from error")
				log.Infof("Waiting for rate limit window")
				time.Sleep(time.Minute)
			}
			rateLimitRetry++
			i = 1
		default:
			break outer
		}
	}

	return resp, err
}

/* error Object structure expected here is as follows:
&struct{
	Payload: &models.ErrorResponse{}
}
here struct is api specific class like AddDeviceDefault, UpdateDeviceDefault.

If struct child variable hierarchy changes, then code needs modification.
*/

func (w *Worker) getRLLimit(lctx *lmctx.LMContext, err error) *types.RateLimits {
	if err == nil {
		return nil
	}
	log := lmlog.Logger(lctx)
	code := util.GetHTTPStatusCodeFromLMSDKError(err)
	if code == http.StatusTooManyRequests {
		XRateLimitLimit := "XRateLimitLimit"
		limit := reflect.ValueOf(err).Elem().FieldByName(XRateLimitLimit)
		if !limit.IsValid() {
			log.Warnf("Field not found:%s", XRateLimitLimit)
			return nil
		}
		XRateLimitWindow := "XRateLimitWindow"
		window := reflect.ValueOf(err).Elem().FieldByName(XRateLimitWindow)
		if !window.IsValid() {
			log.Warnf("Field not found:%s", XRateLimitWindow)
			return nil
		}
		req := &types.RateLimits{
			Limit:  limit.Int(),
			Window: window.Int(),
		}

		return req
	}
	log.Debugf("Rate limit not reached")

	return nil
}

func (w *Worker) setNewLimit(lctx *lmctx.LMContext, req *types.RateLimits, command *types.WorkerCommand) {
	log := lmlog.Logger(lctx)
	// 0.01 is salt to get correct math evaluation
	newLimit := rate.Limit(float64(req.Limit)/float64(req.Window) - 0.01) // nolint: gomnd
	currentLimit, _ := GetCurrentLimit(command.APIInfo.GetPatternKey())
	if newLimit != 0 && currentLimit != newLimit {
		log.Infof("Setting new rate limits for \"%s\": %v", command.APIInfo.GetPatternKey(), newLimit)
		AddLimiter(lctx, command.APIInfo.GetPatternKey(), newLimit)
	}
}
