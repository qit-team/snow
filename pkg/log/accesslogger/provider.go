package accesslogger

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/kernel/container"
	"github.com/hetiansu5/accesslog"
)

const SingletonMain = "access_logger"

type Provider struct{}

func (p *Provider) Register(args ...interface{}) error {
	conf := args[0].(*config.Config)
	instance, err := InitAccessLog(conf.LogHandler, conf.LogDir)
	if err != nil {
		return err
	}
	container.App.SetSingleton(SingletonMain, instance)
	return nil
}

func (p *Provider) Provides() []string {
	return []string{SingletonMain}
}

func (p *Provider) Close() error {
	return nil
}

func GetAccessLogger(args ...string) *accesslog.AccessLogger {
	name := SingletonMain
	if len(args) > 0 {
		if args[0] != "" {
			name = args[0]
		}
	}
	rc, ok := container.App.GetSingleton(name).(*accesslog.AccessLogger)
	if ok {
		return rc
	}
	return nil
}
