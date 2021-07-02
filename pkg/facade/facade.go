package facade

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
)

const responseWaitTimeout = 5 * time.Minute

// Facade implements types.LMFacade interface
type Facade struct {
	mu         sync.RWMutex
	WorkerConf map[int]*types.WConfig
}

// NewFacade creates new facade object
func NewFacade() *Facade {
	f := &Facade{
		WorkerConf: make(map[int]*types.WConfig),
	}

	return f
}

// Send Async command handler
// func (f *Facade) Send(resource string, command types.ICommand) {
//
// }

// SendReceive sync command handler
func (f *Facade) SendReceive(lctx *lmctx.LMContext, command *types.WorkerCommand) (interface{}, error) {
	var res interface{}
	var err error
	res, err = f.sendRecv(lctx, command)

	return res, err
}

// RegisterWorker Registers worker into facade to handler commands of mentioned resource
// plugin pattern, if worker go routine dies for some reason, watcher should create worker and register again
func (f *Facade) RegisterWorker(w types.Worker) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.WorkerConf[w.GetConfig().ID] = w.GetConfig()

	return true, nil
}

// UnregisterWorker Registers worker into facade to handler commands of mentioned resource
// plugin pattern, if worker go routine dies for some reason, watcher should create worker and register again
func (f *Facade) UnregisterWorker(w types.Worker) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.WorkerConf, w.GetConfig().ID)

	return true, nil
}

func (f *Facade) Count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.WorkerConf)
}

func (f *Facade) sendRecv(lctx *lmctx.LMContext, command *types.WorkerCommand) (interface{}, error) {
	log := lmlog.Logger(lctx)
	respch := make(chan *types.WorkerResponse)
	ctx := context.Background()
	command.SetContext(ctx)
	if command.IsSync() {
		log.Debugf("sync command... setting response channel")
		command.SetResponseChannel(respch)
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, responseWaitTimeout)
		defer cancel()
	}

	// Convert to switch case when adding new if conditions
	if err := f.send(lctx, command); err != nil {
		return nil, err
	}
	if command.IsSync() {
		// timeout := time.NewTicker(responseWaitTimeout)
		select {
		case t := <-respch:
			return t.Response, t.Error
		case <-ctx.Done():
			// case <-timeout.C:
			log.Errorf("facade timed out")

			return nil, errors.New("facade timed out for waiting on response")
		}
	}
	return nil, nil
}

func (f *Facade) send(lctx *lmctx.LMContext, command *types.WorkerCommand) error {
	log := lmlog.Logger(lctx)
	partitionKey := lctx.Extract(constants.PartitionKey)
	if partitionKey == nil {
		err := fmt.Errorf("partition_key does not exist in LMContext, cannot send command to worker")
		log.Errorf("%s", err)
		return err
	}

	key, ok := partitionKey.(string)
	if !ok {
		err := fmt.Errorf("partition_key must be string, cannot send command to worker: %v", partitionKey)
		log.Errorf("%s", err)
		return err
	}

	h := fnv.New32a()

	if _, err := h.Write([]byte(key)); err != nil {
		return fmt.Errorf("hash function failed on provided partition key [%s]: %w", key, err)
	}

	var workerID int
	if workers := uint32(len(f.WorkerConf)); workers > 0 {
		workerID = int(h.Sum32() % workers)
	} else {
		return fmt.Errorf("no worker available to execute command [%s]", key)
	}

	log.Debugf("Sending to worker: %d on %v", workerID, f.WorkerConf[workerID].GetChannel())

	f.WorkerConf[workerID].GetChannel() <- command

	return nil
}
