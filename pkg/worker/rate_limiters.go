package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"golang.org/x/time/rate"
)

var (
	rateLimiters = make(map[string]*rate.Limiter)
	mu           sync.Mutex
)

func AddLimiter(lctx *lmctx.LMContext, key string, limit rate.Limit) {
	log := lmlog.Logger(lctx)
	log.Infof("Setting new Limit for %s: %v", key, limit)
	if limiter, ok := getLimiter(key); ok {
		limiter.SetLimit(limit)
	} else {
		rateLimiters[key] = rate.NewLimiter(limit, 1000)
	}
}

func getLimiter(key string) (*rate.Limiter, bool) {
	mu.Lock()
	defer mu.Unlock()
	return rateLimiters[key], rateLimiters[key] != nil
}

func GetCurrentLimit(key string) (rate.Limit, bool) {
	if limiter, ok := getLimiter(key); ok {
		return limiter.Limit(), true
	}
	return rate.Inf, false
}

func Allow(key string, wait bool) (bool, error) {
	if limiter, ok := getLimiter(key); ok {
		switch {
		case limiter.Allow():
			return true, nil
		case wait:
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			if err := limiter.Wait(ctx); err != nil {
				return false, fmt.Errorf("limiter wait failed: %w", err)
			}
		default:
			return false, fmt.Errorf("limiter cannot burst token from limiter of [%s], and wait flag is false", key)
		}
	}
	return false, fmt.Errorf("no limter found with key: %s", key)
}
