package lmlog

import (
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/sirupsen/logrus"
)

// NewLMContextWith creates context with provided log entry
func NewLMContextWith(logger *logrus.Entry) *lmctx.LMContext {
	ctx := lmctx.NewLMContext()
	entryWithDebugID := logger.WithFields(logrus.Fields{"debug_id": util.GetShortUUID()})
	ctx.Set("logger", entryWithDebugID)
	return ctx
}

// Logger returns logger entry from context
func Logger(lctx *lmctx.LMContext) *logrus.Entry {
	return lctx.Extract("logger").(*logrus.Entry)
}

// LMContextWithFields wraps new fields on this context and returns new context
func LMContextWithFields(lctx *lmctx.LMContext, fields logrus.Fields) *lmctx.LMContext {
	entry := Logger(lctx)
	newEntry := entry.WithFields(fields)
	ctx := lmctx.NewLMContext()
	ctx.Set("logger", newEntry)
	return ctx
}
