package healthz

import (
	"net/http"

	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/sirupsen/logrus"
)

// HandleFunc is an http handler function to expose health metrics.
func HandleFunc(w http.ResponseWriter, req *http.Request) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"handler": "healthz"}))
	log := lmlog.Logger(lctx)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Errorf("Failed to write healthz: %v", err)
	}
}
