package main

const (
	_tplConsoleKernel = `package console

import (
	"github.com/robfig/cron"
)

/**
 * 配置执行计划
 * @wiki https://godoc.org/github.com/robfig/cron
 */
func RegisterSchedule(c *cron.Cron) {
	//c.AddFunc("0 30 * * * *", test)
	//c.AddFunc("@hourly", test)
	c.AddFunc("@every 10s", test)
}
`

	_tplConsoleTest = `package console

import "fmt"

func test() {
	fmt.Println("run test")
}
`

	_tplCommand = `package console

import (
	"github.com/qit-team/snow-core/command"
)

func RegisterCommand(c *command.Command) {
	c.AddFunc("test", test)
}
`
)
