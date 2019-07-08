package alimns

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/kernel/container"
	"github.com/aliyun/aliyun-mns-go-sdk"
)

const (
	SingletonMain = "ali_mns"
)

type Provider struct{}

func (p *Provider) Register(args ...interface{}) error {
	conf := args[0].(*config.Config)
	client, err := NewMnsClient(conf.Mns)
	if err != nil {
		return err
	}
	container.App.SetSingleton(SingletonMain, client)
	return nil
}

func (p *Provider) Provides() []string {
	return []string{SingletonMain}
}

//释放资源
func (p *Provider) Close() error {
	arr := p.Provides()
	for _, v := range arr {
		c := GetMns(v)
		if c != nil {
			// Close()
		}
	}
	return nil
}

//获取mns实例
func GetMns(args ...string) ali_mns.MNSClient {
	name := SingletonMain
	if len(args) > 0 {
		if args[0] != "" {
			name = args[0]
		}
	}
	rc, ok := container.App.GetSingleton(name).(ali_mns.MNSClient)
	if ok {
		return rc
	}
	return nil
}
