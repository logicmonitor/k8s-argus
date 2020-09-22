package cronjob

import (
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/robfig/cron/v3"
)

// RegisterCron to register cron job
func RegisterCron(lctx *lmctx.LMContext, cronSpec string, handlerFunc func()) *cron.Cron {
	log := lmlog.Logger(lctx)
	c := cron.New(func(c *cron.Cron) {
		_, err := c.AddFunc(cronSpec, handlerFunc)
		if err != nil {
			log.Errorf("Failed to create cron job. Error: %v", err)
		}
	})
	return c
}
