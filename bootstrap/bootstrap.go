package bootstrap

import (
	"github.com/qit-team/snow-core/db"
	"github.com/qit-team/snow-core/kernel/close"
	"github.com/qit-team/snow-core/kernel/container"
	"github.com/qit-team/snow-core/log/accesslogger"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow/app/jobs"
	"github.com/qit-team/snow/app/jobs/basejob"
	"github.com/qit-team/snow/config"
)

// 全局变量
var App *container.Container

/**
 * 服务引导程序
 */
func Bootstrap(conf *config.Config) (err error) {
	// 容器
	App = container.App

	// 注册db服务
	// 第一个参数为注入别名，第二个参数为配置，第三个参数可选为是否懒加载
	err = db.Pr.Register(db.SingletonMain, conf.Db)
	if err != nil {
		return
	}

	// 注册redis服务
	err = redis.Pr.Register(redis.SingletonMain, conf.Redis)
	if err != nil {
		return
	}

	// 注册mns服务
	// err = alimns.Pr.Register(alimns.SingletonMain, conf.Mns, true)
	// if err != nil {
	// 	return
	// }

	// 注册日志类服务
	err = logger.Pr.Register(logger.SingletonMain, conf.Log, true)
	if err != nil {
		return
	}

	// 注册access log服务
	err = accesslogger.Pr.Register(accesslogger.SingletonMain, conf.Log)
	if err != nil {
		return
	}

	// 注册应用停止时调用的关闭服务
	close.MultiRegister(db.Pr, redis.Pr)

	// 注册job register，为了非job模式的消息入队调用
	basejob.SetJobRegister(jobs.RegisterWorker)
	return nil
}
