package lmlog

import (
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/sirupsen/logrus"
)

// NewLMContextWith creates context with provided log entry
func NewLMContextWith(logger *logrus.Entry) *lmctx.LMContext {
	ctx := lmctx.NewLMContext()
	entryWithDebugID := logger.WithFields(logrus.Fields{"debug_id": getShortUUID()})
	ctx.Set("logger", entryWithDebugID)

	return ctx
}

// Logger returns logger entry from context
func Logger(lctx *lmctx.LMContext) *logrus.Entry {
	return lctx.Extract("logger").(*logrus.Entry).WithField("method", getCallerFunctionName())
}

// LMContextWithFields wraps new fields on this context and returns new context
func LMContextWithFields(lctx *lmctx.LMContext, fields logrus.Fields) *lmctx.LMContext {
	ctx := lctx.Copy()
	entry := Logger(ctx)
	ctx.Set("logger", entry.WithFields(fields).WithFields(logrus.Fields{"p_debug_id": ctx.Extract("debug_id"), "debug_id": getShortUUID()}))

	return ctx
}

func LMContextWithLMResourceID(lctx *lmctx.LMContext, lmid int32) *lmctx.LMContext {
	return LMContextWithFields(lctx, logrus.Fields{"lm_resource_id": lmid})
}

// getShortUUID returns short ids. introduced this util function to start for traceability of events and its logs
func getShortUUID() uint32 {
	return uuid.New().ID()
}

// getCallerFunctionName get caller
func getCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of

	return getNthCallerFunctionName(3) // nolint: gomnd
}

// getNthCallerFunctionName get nth caller
func getNthCallerFunctionName(n int) string {
	// Skip GetCallerFunctionName and the function to get the caller of

	return strings.TrimPrefix(getFrame(n).Function, "github.com/logicmonitor/k8s-argus/")
}

// Referenced here: https://play.golang.org/p/cv-SpkvexuM
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2 // nolint: gomnd

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2) // nolint: gomnd
	frame := runtime.Frame{Function: "unknown"}            // nolint: exhaustivestruct

	n := runtime.Callers(0, programCounters)
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])

		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
