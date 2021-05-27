package resource

import (
	"reflect"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/filters"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/node"
	"github.com/logicmonitor/k8s-argus/pkg/watch/pod"
	"github.com/logicmonitor/k8s-argus/pkg/watch/service"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WatcherConfigurer to provide custom options provider to put extra props on resource
func WatcherConfigurer(resourceType enums.ResourceType) types.WatcherConfigurer {
	switch resourceType {
	case enums.Pods:

		return &pod.Watcher{}
	case enums.Services:

		return &service.Watcher{}
	case enums.Nodes:

		return &node.Watcher{}
	case enums.Deployments, enums.ETCD, enums.Hpas, enums.Namespaces, enums.Unknown:

		return &emptyWatcher{}
	default:

		return &emptyWatcher{}
	}
}

func inferResourceType(newObj interface{}) (enums.ResourceType, bool) {
	// TypeMeta does not exist in object, so we need to infer it explicitly here
	var rt enums.ResourceType
	split := strings.Split(reflect.TypeOf(newObj).String(), ".")
	if err := (&rt).UnmarshalText([]byte(strings.ToLower(split[1]))); err != nil {
		logrus.Errorf("Error while determining ResourceType: [%s] for object [%v]", err, newObj)

		return enums.Unknown, true
	}

	return rt, false
}

func getRootContext(rt enums.ResourceType, newObj interface{}, event string) *lmctx.LMContext {
	objectMeta := rt.ObjectMeta(newObj)
	fields := logrus.Fields{"name": rt.FQName(objectMeta.Name), "type": rt.Singular()}

	if rt.IsNamespaceScopedResource() {
		fields["ns"] = objectMeta.Namespace
	}
	if event != "" {
		fields["event"] = event
	}

	return lmlog.NewLMContextWith(logrus.WithFields(fields))
}

// RecordDeleteEventLatency logs latency of receiving delete event to argus.
func RecordDeleteEventLatency(lctx *lmctx.LMContext, rt enums.ResourceType, obj interface{}) {
	log := lmlog.Logger(lctx)
	if meta := rt.ObjectMeta(obj); meta.DeletionTimestamp != nil {
		// TODO: PROM_METRIC stats: stats of delete event according to object type, time (max, min, average)
		log.Infof("delete event latency %v", time.Since(meta.DeletionTimestamp.Time))
	} else {
		// TODO: PROM_METRIC counter: object without delete timestamp
		log.Warnf("delete event context doesn't have deleteTimestamp on it")
	}
}

// EvaluateResourceExclusion eval
func EvaluateResourceExclusion(lctx *lmctx.LMContext, resourceType enums.ResourceType, meta metav1.ObjectMeta) (bool, error) {
	return filters.Eval(lctx, resourceType, getEvalInput(meta))
}

// getEvalInput generates evaluation parameters based on labels and specified resource
func getEvalInput(meta metav1.ObjectMeta) govaluate.MapParameters {
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
	logrus.Tracef("Eval Input: %v", evaluationParams)

	return evaluationParams
}
