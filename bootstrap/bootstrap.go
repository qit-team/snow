package bootstrap

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/redis"
	"github.com/qit-team/snow/pkg/db"
	"github.com/qit-team/snow/pkg/kernel/container"
	"github.com/qit-team/snow/pkg/kernel/close"
	"github.com/qit-team/snow/pkg/log/accesslogger"
	"github.com/qit-team/snow/pkg/log/logger"
	"github.com/qit-team/snow/pkg/kernel/provider"
	"github.com/qit-team/snow/pkg/alimns"
)

//全局变量
var App *container.Container

/**
 * 服务启动，进行容器注入
 * server包依赖于此函数
*/
func Bootstrap(conf *config.Config) (err error) {
	//容器
	App = container.App

	//db服务
	dbProvider := &db.Provider{}
	//redis服务
	redisProvider := &redis.Provider{}

	//mns服务
	alimnsProvider := &alimns.Provider{}

	//日志类服务
	loggerProvider := &logger.Provider{}

	//access log服务
	accessLoggerProvider := &accesslogger.Provider{}

	err = RegisterProviders(conf, dbProvider, redisProvider, loggerProvider, accessLoggerProvider, alimnsProvider)
	if err != nil {
		return
	}

	//注册应用停止时调用的关闭服务
	close.MultiRegister(dbProvider, redisProvider)
	return nil
}

//批量注册
func RegisterProviders(conf *config.Config, providers ...provider.Provider) error {
	for _, p := range providers {
		err := p.Register(conf)
		if err != nil {
			return err
		}
	}
	return nil
}
