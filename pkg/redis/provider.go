package redis

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/kernel/container"
	redis_pool "github.com/hetiansu5/go-redis-pool"
)

const (
	SingletonMain = "redis"
)

type Provider struct{}

func (p *Provider) Register(args ...interface{}) error {
	conf := args[0].(*config.Config)
	//注入redis
	client, err := NewRedisClient(conf.Redis)
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
		c := GetRedis(v)
		if c != nil {
			c.Close()
		}
	}
	return nil
}

//获取redis实例
func GetRedis(args ...string) *redis_pool.ReplicaPool {
	name := SingletonMain
	if len(args) > 0 {
		if args[0] != "" {
			name = args[0]
		}
	}
	rc, ok := container.App.GetSingleton(name).(*redis_pool.ReplicaPool)
	if ok {
		return rc
	}
	return nil
}
