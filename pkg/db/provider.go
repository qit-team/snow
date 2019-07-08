package db

import (
	"github.com/qit-team/snow/pkg/kernel/container"
	"github.com/qit-team/snow/config"
	"github.com/go-xorm/xorm"
)

const (
	SingletonMain = "db"
)

type Provider struct{}

func (p *Provider) Register(args ...interface{}) error {
	conf := args[0].(*config.Config)
	db, err := NewEngineGroup(conf.Db)
	if err != nil {
		return err
	}
	container.App.SetSingleton(SingletonMain, db)
	return nil
}

func (p *Provider) Provides() []string {
	return []string{SingletonMain}
}

//释放资源
func (p *Provider) Close() error {
	arr := p.Provides()
	for _, v := range arr {
		c := GetDb(v)
		if c != nil {
			c.Close()
		}
	}
	return nil
}

func GetDb(args ...string) (*xorm.EngineGroup) {
	name := SingletonMain
	if len(args) > 0 {
		if args[0] != "" {
			name = args[0]
		}
	}
	rc, ok := container.App.GetSingleton(name).(*xorm.EngineGroup)
	if ok {
		return rc
	}
	return nil
}
