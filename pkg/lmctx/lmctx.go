package lmctx

import (
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
)

// LMContext will be used to pass metadata from one to other, for ex: log context
type LMContext struct {
	kv map[string]interface{}
}

// Extract will return value against key
func (lc *LMContext) Extract(key string) interface{} {
	return lc.kv[key]
}

// Set will store key val in map
func (lc *LMContext) Set(key string, val interface{}) {
	lc.kv[key] = val
}

// Logger returns logger entry from context
func (lc *LMContext) Logger() *logrus.Entry {
	return lc.Extract("logger").(*logrus.Entry)
}

// WithFields wraps new fields on this context and returns new context
func (lc *LMContext) WithFields(fields logrus.Fields) *LMContext {
	entry := lc.Logger()
	newEntry := entry.WithFields(fields)
	return WithLogger(newEntry)
}

// NewLMContext creates new context object
func NewLMContext() *LMContext {
	ctx := &LMContext{
		kv: make(map[string]interface{}),
	}
	return ctx
}

// WithLogger creates context with provided log entry
func WithLogger(logEntry *logrus.Entry) *LMContext {
	ctx := NewLMContext()
	entryWithDebugID := logEntry.WithFields(logrus.Fields{"debug_id": util.GetShortUUID()})
	ctx.Set("logger", entryWithDebugID)
	return ctx
}
