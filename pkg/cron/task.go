package cron

import "vid-msa/pkg/logger"

// testfun 定时任务测试函数
func testfun() {
	logger.Info("task1 is running")
}
