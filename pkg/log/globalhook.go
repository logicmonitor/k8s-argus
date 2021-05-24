package lmlog

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// DefaultFieldHook global hook to add pod id in log stmt
type DefaultFieldHook struct{}

// Fire log hook method before dump
func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["argus_pod_id"] = os.Getenv("MY_POD_NAME")
	entry.Data["goroutine"] = goID()

	return nil
}

// Levels level array to enable hook on
func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func goID() int {
	var id int
	var err error
	defer func() {
		if r := recover(); r != nil {
			id = -1
		}
	}()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err = strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return id
}
