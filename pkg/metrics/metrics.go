package metrics

import (
	"encoding/json"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	m    *metrics
	once sync.Once
)

type metrics struct {
	APIErrors  uint
	RESTErrors uint
}

func init() {
	once.Do(func() {
		m = &metrics{
			APIErrors:  0,
			RESTErrors: 0,
		}
	})
}

// APIError increments the API error count by 1.
func APIError() {
	m.APIErrors++
}

// RESTError increments the REST error count by 1.
func RESTError() {
	m.RESTErrors++
}

// HandleFunc is an http handler function to expose Argus metrics.
func HandleFunc(w http.ResponseWriter, req *http.Request) {
	resp, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Errorf("Failed to write metrics: %v", err)
	}
}
