package commands

import (
	"github.com/qit-team/snow-core/command"
	"fmt"
)

func RegisterCommand(c *command.Command) {
	c.AddFunc("test", test)
}

func test() {
	fmt.Println("run test command")
	return
}
