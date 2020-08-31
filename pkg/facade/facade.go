package facade

import (
	"errors"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/sirupsen/logrus"
)

// Facade implements types.LMFacade interface
type Facade struct {
	WorkerConf map[string]*types.WConfig
}

// NewFacade creates new facade object
func NewFacade() *Facade {
	f := &Facade{
		WorkerConf: make(map[string]*types.WConfig),
	}
	return f
}

// Send Async command handler
//func (f *Facade) Send(resource string, command types.ICommand) {
//
//}

// SendReceive sync command handler
func (f *Facade) SendReceive(lctx *lmctx.LMContext, resource string, command types.ICommand) (interface{}, error) {
	var res interface{}
	var err error
	res, err = f.sendRecv(lctx, resource, command)
	return res, err
}

// RegisterWorker Registers worker into facade to handler commands of mentioned resource
// plugin pattern, if worker go routine dies for some reason, watcher should create worker and register again
func (f *Facade) RegisterWorker(resource string, w types.Worker) (bool, error) {
	logrus.Debugf("registering worker for %s %#v", resource, w.GetConfig())
	f.WorkerConf[resource] = w.GetConfig()
	return true, nil
}

func (f *Facade) sendRecv(lctx *lmctx.LMContext, resource string, command types.ICommand) (interface{}, error) {
	log := lmlog.Logger(lctx)
	respch := make(chan *types.WorkerResponse)
	var i interface{} = command
	if cmd, ok := i.(types.Responder); ok {
		log.Debugf("sync command... setting response channel")
		cmd.SetResponseChannel(respch)
	}

	switch i.(type) {
	case types.IHTTPCommand:
		wc := f.WorkerConf[resource]
		ch := wc.GetChannel(command)
		ch <- command
		timeout := time.NewTicker(5 * time.Minute)
		select {
		case t := <-respch:
			return t.Response, t.Error
		case <-timeout.C:
			log.Errorf("facade timed out")
			return nil, errors.New("facade timed out for waiting on response")
		}
	}

	return nil, errors.New("unknown Command")

}
