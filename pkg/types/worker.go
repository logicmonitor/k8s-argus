package types

import (
	"context"
	"fmt"
	"net/http"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
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
	return wc.ExecFunc()
}
