package eventprocessor

import "github.com/logicmonitor/k8s-argus/pkg/lmctx"

// RunnerRequest function type to point to execute function
type RunnerRequest func()

type RunnerCommand struct {
	ExecFunc RunnerRequest
	Lctx     *lmctx.LMContext
}

func (rc *RunnerCommand) Execute() {
	rc.ExecFunc()
}

// RunnerConfig runner configuration
type RunnerConfig struct {
	ID   int
	inCh chan *RunnerCommand
}

func (rc *RunnerConfig) GetChannel() chan *RunnerCommand {
	return rc.inCh
}

// NewRunnerConfig runner configuration
func NewRunnerConfig(id int, size int) *RunnerConfig {
	ch := make(chan *RunnerCommand, size)

	return &RunnerConfig{
		inCh: ch,
		ID:   id,
	}
}

// EventRunner worker interface to provide interface method
type EventRunner interface {
	Run()
	GetConfig() *RunnerConfig
}

// RunnerFacade public interface others to interact with
type RunnerFacade interface {
	// Async
	// Send(command ICommand)

	// Send sync api call
	Send(*lmctx.LMContext, *RunnerCommand) error
	// RegisterRunner registers worker to facade client to put command objects on channel
	RegisterRunner(EventRunner) (bool, error)
	// UnregisterRunner registers worker to facade client to put command objects on channel
	UnregisterRunner(EventRunner) (bool, error)
	// Count registers worker to facade client to put command objects on channel
	Count() int
}
