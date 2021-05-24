package worker

import (
	"context"
	"errors"
	"time"
)

const millisecondsOfMinute = 60000

// RLTokenizer tokenizer which contains regulated token generation
type RLTokenizer struct {
	ch     <-chan interface{}
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRLTokenizer creates new tokenizer with mentioned limit and starts it implicitly
func NewRLTokenizer(limit int) *RLTokenizer {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan interface{}, limit)

	go func(wch chan<- interface{}) {
		ticker := time.NewTicker(time.Duration(millisecondsOfMinute/limit) * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				close(wch)
				return
			case <-ticker.C:
				wch <- 1
			}
		}
	}(ch)

	return &RLTokenizer{
		ch:     ch,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (rlt *RLTokenizer) popToken() error {
	if rlt.ctx.Err() != nil {
		return rlt.ctx.Err()
	}
	select {
	case <-rlt.ch:

		return nil
	case <-time.After(1 * time.Minute):

		return errors.New("new token did not received in 1 minute, reconfigure tokenizer")
	}
}

func (rlt *RLTokenizer) stop() {
	rlt.cancel()
}
