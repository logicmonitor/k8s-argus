package worker

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	ratelimiter "github.com/logicmonitor/k8s-argus/pkg/rl"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
)

const (
	// retryLimit Default retry limit for workers
	retryLimit = 2
	// maxRateLimitRetry number of retries to give up on rate limit retry request
	maxRateLimitRetry = 2
	// retryBackoffTimeDurationDefault default
	retryBackoffTimeDurationDefault = 5 * time.Second
	// sendResponseTimeout send response on channel timeout
	sendResponseTimeout = 2 * time.Millisecond
	// workerIdleNotifyTimeout worker idle timeout
	workerIdleNotifyTimeout = 10 * time.Second
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
func NewHTTPWorker(c *types.WConfig) *Worker {
	w := &Worker{
		config:         c,
		initialized:    false,
		running:        false,
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
	for _, ch := range w.config.Channels {
		channelMap[ch] = true
	}
	for ch := range channelMap {
		go func(inch chan types.ICommand) {
			defer func() {
				w.running = false
			}()
			for {
				timeout := time.NewTicker(workerIdleNotifyTimeout)
				select {
				case command := <-inch:
					lctx2 := command.LMContext()
					commandCtx := lmlog.LMContextWithFields(lctx2, logrus.Fields{"worker": w.config.ID})
					log2 := lmlog.Logger(commandCtx)
					log2.Debugf("Received command")
					w.handleCommand(commandCtx, command)
				case <-timeout.C:
					log.Tracef("%v worker is idle", w.config.ID)
				}
			}
		}(ch)
	}
	uch := make(chan types.WorkerRateLimitsUpdate)
	w.watchForLimitChanges(lctx, uch)
	ratelimiter.RegisterWorkerNotifyChannel(w.config.ID, uch)
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
	// Convert to switch case when adding new if conditions
	if cmdRef, ok := tmpCommand.(types.IHTTPCommand); ok {
		tch := w.getTokenizer(cmdRef.GetCategory(), cmdRef.GetMethod())
		err := tch.popToken()
		if err != nil && errors.Is(err, context.Canceled) {
			tch = w.getTokenizer(cmdRef.GetCategory(), cmdRef.GetMethod())
			tch.popToken() // nolint: errcheck, gosec
		}
	}
}

func (w *Worker) handleCommand(lctx *lmctx.LMContext, command types.ICommand) {
	log := lmlog.Logger(lctx)
	log.Tracef("Poping token")
	w.popRLToken(command)
	log.Tracef("Token popped")
	retryMax := w.config.RetryLimit
	if retryMax < 1 {
		retryMax = retryLimit
	}
	resp, err := w.executeWithRetry(lctx, retryMax, command)
	wresp := &types.WorkerResponse{
		Response: resp,
		Error:    err,
	}
	if n, ok := command.(types.Responder); ok {
		log.Tracef("Sync command")
		if rch := n.GetResponseChannel(); rch != nil {
			log.Tracef("Response channel defined %v", wresp)
			select {
			case rch <- wresp:
			// to make non blocking, golang unbuffered channels are blocking if there is no goroutine listening on channel. if facade timed out and returned then this would become blocking
			case <-time.After(sendResponseTimeout):
				log.Warnf("Response cannot sent")
			}
		}
	}
}

func getHTTPStatusCode(err error) int {
	errRegex := regexp.MustCompile(`(?P<api>\[.*\])\[(?P<code>\d+)\].*`)
	matches := errRegex.FindStringSubmatch(err.Error())
	if len(matches) < 3 { // nolint: gomnd
		return -1
	}

	code, err := strconv.Atoi(matches[2])
	if err != nil {
		return -1
	}

	return code
}

func (w *Worker) executeWithRetry(lctx *lmctx.LMContext, retry int, command types.ICommand) (interface{}, error) {
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
		code := getHTTPStatusCode(err)
		log.Debugf("Status code: %v", code)

		// retry only if request failed because of server error
		switch {
		// 5XX server error
		case code/100 == 5: // nolint: gomnd
			log.Warningf("Request failed with error %v, retrying for %v time...", err, i)
			time.Sleep(retryBackoffTimeDurationDefault)
		case code == http.StatusTooManyRequests:
			req := w.getRLLimit(lctx, command, err)
			w.sendRLLimitUpdate(lctx, req)
			log.Infof("Waiting for rate limit window %v", req.Window)
			time.Sleep(time.Duration(req.Window) * time.Second)
			rateLimitRetry++
			i = 1
		default:
			break outer
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
	// Convert to switch case while adding new if conditions
	if cmdRef, ok := tmpCommand.(types.LMHCErrParser); ok {
		hc := command.(types.IHTTPCommand)
		errResp := cmdRef.ParseErrResponse(err)
		code := getHTTPStatusCode(err)
		if code == http.StatusTooManyRequests {
			headers := errResp.ErrorDetail.(map[string]interface{}) // nolint: forcetypeassert
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
		log.Debugf("Sending update to ratelimiter %v", req)
		ratelimiter.GetUpdateRequestChannel() <- *req
	}
}
