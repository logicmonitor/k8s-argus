package types

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
)

type HTTPMethod string

const (
	// HTTPGet get
	HTTPGet = HTTPMethod(http.MethodGet)

	// HTTPDelete delete
	HTTPDelete = HTTPMethod(http.MethodDelete)

	// HTTPPatch patch
	HTTPPatch = HTTPMethod(http.MethodPatch)

	// HTTPPost post
	HTTPPost = HTTPMethod(http.MethodPost)

	// HTTPPut put
	HTTPPut = HTTPMethod(http.MethodPut)
)

type APIInfo struct {
	URLPattern string
	Method     HTTPMethod
}

func (apiInfo *APIInfo) GetPatternKey() string {
	return fmt.Sprintf("%s;%s", apiInfo.URLPattern, apiInfo.Method)
}

type WorkerCommand struct {
	ExecFunc     ExecRequest
	Lctx         *lmctx.LMContext
	ctx          context.Context
	APIInfo      APIInfo
	isAsync      bool
	responseChan chan *WorkerResponse
}

func (wc *WorkerCommand) IsSync() bool {
	return !wc.isAsync
}

func (wc *WorkerCommand) SetResponseChannel(ch chan *WorkerResponse) {
	wc.responseChan = ch
}

func (wc *WorkerCommand) GetResponseChannel() chan *WorkerResponse {
	return wc.responseChan
}

func (wc *WorkerCommand) SetContext(ctx context.Context) {
	wc.ctx = ctx
}

func (wc *WorkerCommand) GetContext() context.Context {
	return wc.ctx
}

func (wc *WorkerCommand) Execute() (interface{}, error) {
	start := time.Now()
	code := http.StatusOK
	defer func() {
		metrics.APIResponseTimeSummary.
			WithLabelValues(wc.APIInfo.URLPattern,
				string(wc.APIInfo.Method),
				fmt.Sprintf("%d", code)).
			Observe(float64(time.Since(start).Nanoseconds()))
		// TODO: following needs to remove when openmetrics collection works with aggregate
		metrics.APIResponseTimeSummary.
			WithLabelValues(wc.APIInfo.URLPattern,
				string(wc.APIInfo.Method),
				fmt.Sprintf("%d%s", code/100, "XX")). // nolint: gomnd
			Observe(float64(time.Since(start).Nanoseconds()))
		// TODO: following needs to remove when openmetrics collection works with aggregate
		metrics.APIResponseTimeSummary.
			WithLabelValues(wc.APIInfo.URLPattern, "all", fmt.Sprintf("%d", code)).
			Observe(float64(time.Since(start).Nanoseconds()))
	}()
	resp, err := wc.ExecFunc()
	code = GetCode(resp, err)
	return resp, err
}

func GetCode(resp interface{}, err error) int {
	code := http.StatusOK
	if err != nil {
		code = GetHTTPStatusCodeFromLMSDKError(err)
	} else {
		e2, ok := resp.(error)
		if ok {
			if c := GetHTTPStatusCodeFromLMSDKError(e2); c > 0 {
				code = c
			}
		}
	}
	return code
}

// GetHTTPStatusCodeFromLMSDKError get code
func GetHTTPStatusCodeFromLMSDKError(err error) int {
	if err == nil {
		return -2
	}
	if errors.Is(err, context.DeadlineExceeded) {
		// 408 client timeout error
		return http.StatusRequestTimeout
	}
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
