package eventprocessor

import (
	"fmt"
	"hash/fnv"
	"sync"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
)

// RFacade implements types.RunnerFacade interface
type RFacade struct {
	mu         sync.RWMutex
	RunnerConf map[int]*RunnerConfig
}

// NewRFacade creates new facade object
func NewRFacade() *RFacade {
	f := &RFacade{
		RunnerConf: make(map[int]*RunnerConfig),
	}

	return f
}

func (rf *RFacade) Send(lctx *lmctx.LMContext, runnerCommand *RunnerCommand) error {
	log := lmlog.Logger(lctx)
	partitionKey := lctx.Extract(constants.PartitionKey)
	if partitionKey == nil {
		err := fmt.Errorf("partition_key does not exist in LMContext, cannot send command to runner")
		log.Errorf("%s", err)
		return err
	}

	key, ok := partitionKey.(string)
	if !ok {
		err := fmt.Errorf("partition_key must be string, cannot send command to runner: %v", partitionKey)
		log.Errorf("%s", err)
		return err
	}

	h := fnv.New32a()

	if _, err := h.Write([]byte(key)); err != nil {
		return fmt.Errorf("hash function failed on provided partition key [%s]: %w", key, err)
	}

	var workerID int
	if workers := uint32(len(rf.RunnerConf)); workers > 0 {
		workerID = int(h.Sum32() % workers)
	} else {
		return fmt.Errorf("no runner available to execute command [%s]", key)
	}

	log.Debugf("Sending to runner: %d on %v", workerID, rf.RunnerConf[workerID].GetChannel())

	rf.RunnerConf[workerID].GetChannel() <- runnerCommand

	return nil
}

func (rf *RFacade) RegisterRunner(eventRunner EventRunner) (bool, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.RunnerConf[eventRunner.GetConfig().ID] = eventRunner.GetConfig()

	return true, nil
}

func (rf *RFacade) UnregisterRunner(eventRunner EventRunner) (bool, error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	delete(rf.RunnerConf, eventRunner.GetConfig().ID)

	return true, nil
}

func (rf *RFacade) Count() int {
	rf.mu.RLock()
	defer rf.mu.RUnlock()
	return len(rf.RunnerConf)
}
