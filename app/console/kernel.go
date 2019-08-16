package console

import (
	"github.com/robfig/cron"
)

/**
 * 配置执行计划
 * @wiki https://godoc.org/github.com/robfig/cron
 */
func RegisterSchedule(c *cron.Cron) {
	// c.AddFunc("0 30 * * * *", test)
	// c.AddFunc("@hourly", test)
	c.AddFunc("@every 10s", test)
}
