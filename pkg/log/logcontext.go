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
	return lctx.Extract("logger").(*logrus.Entry).WithField("method", util.GetCallerFunctionName())
}

// LMContextWithFields wraps new fields on this context and returns new context
func LMContextWithFields(lctx *lmctx.LMContext, fields logrus.Fields) *lmctx.LMContext {
	entry := Logger(lctx)
	newEntry := entry.WithFields(fields)
	ctx := lctx.Copy()
	ctx.Set("logger", newEntry)

	return ctx
}

func LMContextWithLMResourceID(lctx *lmctx.LMContext, lmid int32) *lmctx.LMContext {
	return LMContextWithFields(lctx, logrus.Fields{"lm_resource_id": lmid})
}
