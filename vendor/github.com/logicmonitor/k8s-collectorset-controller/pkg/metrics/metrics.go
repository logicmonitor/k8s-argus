package metrics

import (
	"expvar"
	"runtime"
	"sync"
)

var (
	m    *expvar.Map
	once sync.Once
)

func init() {
	once.Do(func() {
		m = expvar.NewMap("errors")
		m.Add("APIErrors", 0)
		m.Add("RESTErrors", 0)
	})
	expvar.Publish("goroutines", expvar.Func(goroutines))
}

// APIError increments the API error count by 1.
func APIError() {
	m.Add("APIErrors", 1)
}

// RESTError increments the REST error count by 1.
func RESTError() {
	m.Add("RESTErrors", 1)
}

func goroutines() interface{} {
	return runtime.NumGoroutine()
}
