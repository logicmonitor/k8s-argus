package worker

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	rlm "github.com/logicmonitor/k8s-argus/pkg/rl"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
)

const (
	// RETRYLIMIT Default retry limit for workers
	RETRYLIMIT = 2
	// MaxRateLimitRetry number of retries to give up on rate limit retry request
	MaxRateLimitRetry = 2
)

// Worker object
type Worker struct {
	config         *types.WConfig
	initialized    bool
	running        bool
	tokenizers     *sync.Map
	cancelContexts map[<-chan interface{}]func()
}

// GetConfig returns the config object
func (w *Worker) GetConfig() *types.WConfig {
	return w.config
}

// NewWorker creates a worker with provided config
func NewWorker(c *types.WConfig) *Worker {
	return &Worker{
		config:         c,
		initialized:    false,
		running:        false,
		tokenizers:     &sync.Map{},
		cancelContexts: make(map[<-chan interface{}]func()),
	}
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
		tokenizers:     &sync.Map{},
		cancelContexts: make(map[<-chan interface{}]func()),
	}
	return w
}

// Run creates go routine to handle requests
func (w *Worker) Run() {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"worker": w.config.ID}))
	log := lmlog.Logger(lctx)
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
					commandCtx := lmlog.LMContextWithFields(lctx, logrus.Fields{"worker": w.config.ID})
					log = lmlog.Logger(commandCtx)
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
	uch := make(chan types.WorkerRateLimitsUpdate)
	w.watchForLimitChanges(lctx, uch)
	rlm.RegisterWorkerNotifyChannel(w.config.ID, uch)
	w.initialized = true
	w.running = true
}

func (w *Worker) createTokenizer(key string, limit int) *RLTokenizer {
	oldChan, ok := w.tokenizers.Load(key)
	if ok {
		oldChan.(*RLTokenizer).stop()
	}
	newTokenizer := NewRLTokenizer(limit)
	w.tokenizers.Store(key, newTokenizer)
	return newTokenizer
}

func (w *Worker) updateTokenizer(up types.WorkerRateLimitsUpdate) {
	w.createTokenizer(up.Category+up.Method, int(up.Limit))
}
func (w *Worker) watchForLimitChanges(lctx *lmctx.LMContext, ch <-chan types.WorkerRateLimitsUpdate) {
	log := lmlog.Logger(lctx)
	go func(uchan <-chan types.WorkerRateLimitsUpdate) {
		for {
			update := <-uchan
			log.Debugf("Updated received to worker %s: %v", w.config.ID, update)
			w.updateTokenizer(update)
		}
	}(ch)
}

func (w *Worker) getTokenizer(category string, method string) *RLTokenizer {
	ch, ok := w.tokenizers.Load(category + method)
	if !ok {
		return w.createTokenizer(category+method, 1000)
	}
	return ch.(*RLTokenizer)
}

func (w *Worker) popRLToken(command types.ICommand) {
	var tmpCommand interface{} = command
	switch cmdRef := tmpCommand.(type) {
	case types.IHTTPCommand:
		tch := w.getTokenizer(cmdRef.GetCategory(), cmdRef.GetMethod())
		err := tch.popToken()
		if err != nil && err == context.Canceled {
			tch = w.getTokenizer(cmdRef.GetCategory(), cmdRef.GetMethod())
			tch.popToken() // nolint: errcheck, gosec
		}
	}
}

func (w *Worker) handleCommand(lctx *lmctx.LMContext, command types.ICommand) {
	log := lmlog.Logger(lctx)
	log.Debugf("Poping token")
	w.popRLToken(command)
	log.Debugf("Token popped")
	retryMax := w.config.RetryLimit
	if retryMax < 1 {
		retryMax = RETRYLIMIT
	}
	resp, err := w.executeWithRetry(lctx, retryMax, command)
	wresp := &types.WorkerResponse{
		Response: resp,
		Error:    err,
	}
	if n, ok := command.(types.Responder); ok {
		log.Debugf("Sync command")
		if rch := n.GetResponseChannel(); rch != nil {
			log.Debugf("Response channel defined %v", wresp)
			select {
			case rch <- wresp:
			// to make non blocking, golang unbuffered channels are blocking if there is no goroutine listening on channel. if facade timed out and returned then this would become blocking
			case <-time.After(2 * time.Millisecond):
				log.Warnf("Response cannot sent")
			}
		}
	}
}

func (w *Worker) executeWithRetry(lctx *lmctx.LMContext, retry int, command types.ICommand) (interface{}, error) {
	log := lmlog.Logger(lctx)
	var resp interface{}
	var err error
	rateLimitRetry := 0
	for i := 1; i <= retry && rateLimitRetry < MaxRateLimitRetry; i++ {
		resp, err = command.Execute()
		if err == nil {
			break
		}
		code := util.GetHTTPStatusCodeFromLMSDKError(err)
		log.Debugf("Status code: %v", code)

		// retry only if request failed because of server error
		if code >= 500 && code <= 599 {
			log.Warningf("Request failed with error %v, retrying for %v time...", err, i)
			time.Sleep(5 * time.Second)
		} else if code == http.StatusTooManyRequests { // rate limit reached then send update to rate limit manager
			req := w.getRLLimit(lctx, command, err)
			w.sendRLLimitUpdate(lctx, req)
			log.Infof("Waiting for rate limit window %v", req.Window)
			time.Sleep(time.Duration(req.Window) * time.Second)
			//time.Sleep(5 * time.Second)
			// reset looper so that it will retry otherwise if i=retry then it will come out of loop
			rateLimitRetry++
			i = 1
		} else {
			break
		}
	}
	return resp, err
}

func (w *Worker) getRLLimit(lctx *lmctx.LMContext, command types.ICommand, err error) *types.RateLimitUpdateRequest {
	if err == nil {
		return nil
	}
	log := lmlog.Logger(lctx)
	var tmpCommand interface{} = command
	switch cmdRef := tmpCommand.(type) {
	case types.LMHCErrParser:
		hc := command.(types.IHTTPCommand)
		errResp := cmdRef.ParseErrResponse(err)
		code := util.GetHTTPStatusCodeFromLMSDKError(err)
		if code == http.StatusTooManyRequests {
			headers := errResp.ErrorDetail.(map[string]interface{})
			limit, lerr := strconv.Atoi(headers["x-rate-limit-limit"].(string))
			if lerr != nil {
				return nil
			}
			window, werr := strconv.Atoi(headers["x-rate-limit-window"].(string))
			if werr != nil {
				return nil
			}
			req := &types.RateLimitUpdateRequest{
				Worker:   w.config.ID,
				Category: hc.GetCategory(),
				Method:   hc.GetMethod(),
				Limit:    int64(limit),
				Window:   window,
			}
			return req
		}
		log.Debugf("Rate limit not reached")
	}
	return nil
}

func (w *Worker) sendRLLimitUpdate(lctx *lmctx.LMContext, req *types.RateLimitUpdateRequest) {
	log := lmlog.Logger(lctx)
	if req != nil {
		log.Debugf("Sending update to rlm %v", req)
		rlm.GetUpdateRequestChannel() <- *req
	}
}
