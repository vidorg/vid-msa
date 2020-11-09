package cron

import (
	"github.com/robfig/cron/v3"
)

// Init 初始化定时任务
func StartTasks() error {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0/3 * * * * ? ", testfun)
	if err != nil {
		return err
	}
	c.Start()
	return nil
}
