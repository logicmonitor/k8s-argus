package healthz

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleFunc is an http handler function to expose health metrics.
func HandleFunc(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Errorf("Failed to write healthz: %v", err)
	}
}
