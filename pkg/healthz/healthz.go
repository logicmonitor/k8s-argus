package healthz

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// HandleFunc is an http handler function to expose health metrics.
func HandleFunc(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		logrus.Errorf("Failed to write healthz: %v", err)
	}
}
