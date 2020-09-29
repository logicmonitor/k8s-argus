package cronjob

import (
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/robfig/cron/v3"
)

var (
	cj *cron.Cron
)

func init() {
	cj = cron.New()
	// Start the cron scheduler in its own goroutine, or no-op if already started.
	cj.Start()
}

// RegisterFunc adds a func to the Cron to be run on the given schedule.
func RegisterFunc(lctx *lmctx.LMContext, cronSpec string, handlerFunc func()) cron.EntryID {
	log := lmlog.Logger(lctx)
	entryID, err := cj.AddFunc(cronSpec, handlerFunc)
	if err != nil {
		log.Errorf("Failed to add a func to the cron. Error: %v", err)
		return 0
	}
	return entryID
}
