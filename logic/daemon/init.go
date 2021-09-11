package daemon

import (
	"context"
	"github.com/robfig/cron/v3"
)

var cronClient *cron.Cron

func InitCron() {
	cronClient = cron.New(cron.WithSeconds())

	cronClient.Start()
}

func StopCron(ctx context.Context) {
	if cronClient == nil {
		return
	}

	select {
	case <-cronClient.Stop().Done():
		return
	case <-ctx.Done():
		return
	}
}
