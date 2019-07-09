package logger

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/kernel/container"
	"github.com/sirupsen/logrus"
	"os"
)

const SingletonMain = "logger"

type Provider struct{}

func (p *Provider) Register(args ...interface{}) error {
	conf := args[0].(*config.Config)
	instance, err := InitLog(conf.LogHandler, conf.LogDir, conf.LogLevel)
	if err != nil {
		return err
	}
	container.App.SetSingleton(SingletonMain, instance)
	return nil
}

func (p *Provider) Provides() []string {
	return []string{SingletonMain}
}

//释放资源
func (p *Provider) Close() error {
	log, ok := GetLogger().Out.(*os.File)
	if ok {
		log.Sync()
		log.Close()
	}
	return nil
}

func GetLogger(args ...string) *logrus.Logger {
	name := SingletonMain
	if len(args) > 0 {
		if args[0] != "" {
			name = args[0]
		}
	}
	rc, ok := container.App.GetSingleton(name).(*logrus.Logger)
	if ok {
		return rc
	}
	return nil
}
