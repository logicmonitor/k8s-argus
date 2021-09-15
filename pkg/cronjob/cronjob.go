package cronjob

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var cj *cron.Cron

func Init() {
	cj = cron.New()
	// Start the cron scheduler in its own goroutine, or no-op if already started.
	cj.Start()
}

// RegisterFunc adds a func to the Cron to be run on the given schedule.
func RegisterFunc(cronSpec string, handlerFunc func()) (cron.EntryID, error) {
	entryID, err := cj.AddFunc(cronSpec, handlerFunc)
	if err != nil {
		return 0, fmt.Errorf("failed to add cron job: %w", err)
	}

	return entryID, nil
}
