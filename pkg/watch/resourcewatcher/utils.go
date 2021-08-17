package resourcewatcher

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/eventprocessor"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/metrics"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WatcherConfigurer to provide custom options provider to put extra props on resource
func WatcherConfigurer(resourceType enums.ResourceType) types.WatcherConfigurer {
	// nolint: exhaustive
	switch resourceType {
	case enums.Pods:

		return &pod.Watcher{}
	case enums.Services:

		return &service.Watcher{}
	case enums.Nodes:

		return &node.Watcher{}
	default:

		return &emptyWatcher{}
	}
}

func InferResourceType(newObj interface{}) (enums.ResourceType, bool) {
	// TypeMeta does not exist in object, so we need to infer it explicitly here
	var rt enums.ResourceType
	objectTypeStr := reflect.TypeOf(newObj).String()
	split := strings.Split(objectTypeStr, ".")
	var objectName string
	if len(split) < 2 { // nolint: gomnd
		objectName = split[0]
	} else {
		objectName = split[1]
	}

	if err := (&rt).UnmarshalText([]byte(strings.ToLower(split[1]))); err != nil {
		logrus.Errorf("\"%s\" is not valid ResourceType, error: [%s] for object [%v]", objectName, err, newObj)

		return enums.Unknown, false
	}

	return rt, true
}

func getRootContext(lctx *lmctx.LMContext, rt enums.ResourceType, newObj interface{}, event string) *lmctx.LMContext {
	objectMeta := rt.ObjectMeta(newObj)
	fields := logrus.Fields{"name": rt.FQName(objectMeta.Name), "type": rt.Singular()}

	if rt.IsNamespaceScopedResource() {
		fields["ns"] = objectMeta.Namespace
	}
	if event != "" {
		fields["event"] = event
	}
	if conf, err := config.GetConfig(lctx); err == nil {
		fields["display_name"] = util.GetDisplayName(rt, objectMeta, conf)
	}

	clctx := lmlog.LMContextWithFields(lctx, fields)
	clctx.Set(constants.PartitionKey, fmt.Sprintf("%s-%s", rt.String(), objectMeta.Name))
	return clctx
}

// RecordDeleteEventLatency logs latency of receiving delete event to argus.
func RecordDeleteEventLatency(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	log := lmlog.Logger(lctx)
	if meta := rt.ObjectMeta(obj); meta.DeletionTimestamp != nil {
		metrics.DeleteEventLatencySummary.WithLabelValues(rt.String()).Observe(float64(time.Since(meta.DeletionTimestamp.Time).Nanoseconds()))
		// TODO: PROM_METRIC stats: stats of delete event according to object type, time (max, min, average)
		log.Infof("delete event latency %v", time.Since(meta.DeletionTimestamp.Time))
	} else {
		metrics.DeleteEventMissingTimestamp.WithLabelValues(rt.String()).Inc()
		// TODO: PROM_METRIC counter: object without delete timestamp
		log.Warnf("delete event context doesn't have deleteTimestamp on it")
	}
}

// EvaluateResourceExclusion eval
func EvaluateResourceExclusion(lctx *lmctx.LMContext, resourceType enums.ResourceType, meta *metav1.PartialObjectMetadata) (bool, error) {
	return filters.Eval(lctx, resourceType, getEvalInput(lctx, meta))
}

// getEvalInput generates evaluation parameters based on labels and specified resource
func getEvalInput(lctx *lmctx.LMContext, meta *metav1.PartialObjectMetadata) govaluate.MapParameters {
	log := lmlog.Logger(lctx)
	evaluationParams := make(govaluate.MapParameters)
	// adding annotations first and then labels, so that labels get higher precedence
	for key, value := range meta.Annotations {
		key = filters.SanitiseEvalInput(key)
		value = filters.SanitiseEvalInput(value)
		evaluationParams[key] = value
	}
	for key, value := range meta.Labels {
		key = filters.SanitiseEvalInput(key)
		value = filters.SanitiseEvalInput(value)
		evaluationParams[key] = value
	}

	evaluationParams["name"] = filters.SanitiseEvalInput(meta.Name)
	evaluationParams["ns"] = filters.SanitiseEvalInput(meta.Namespace)
	evaluationParams["namespace"] = filters.SanitiseEvalInput(meta.Namespace)
	log.Tracef("Eval Input: %v", evaluationParams)

	return evaluationParams
}

func sendToFacade(facade eventprocessor.RunnerFacade, lctx *lmctx.LMContext, rt enums.ResourceType, event string, function func()) {
	log := lmlog.Logger(lctx)
	defer metrics.ObserveTime(metrics.StartTimeObserver(metrics.ResourceHandlerProcessingTimeSummary.WithLabelValues(rt.String(), event)))
	if err := facade.Send(lctx, &eventprocessor.RunnerCommand{
		ExecFunc: function,
		Lctx:     lctx,
	}); err != nil {
		log.Errorf("Failed to perform event processing: %s", err)

		return
	}
	log.Debug("event processing completed")
}
