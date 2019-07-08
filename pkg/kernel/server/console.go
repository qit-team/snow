package server

import (
	"github.com/qit-team/snow/config"
	"github.com/robfig/cron"
	"fmt"
	"time"
)

func waitConsoleStop(c *cron.Cron) {
	waitStop()

	//暂停新的Cron任务执行
	c.Stop()

	//等待执行中的cron任务结束，目前简单实现等待5s后结束
	fmt.Println("wait 5 sencods")
	time.Sleep(time.Second * 5)
}

// Start Cron Schedule
func StartConsole(confFile, pidFile string, boot func(*config.Config) error, registerSchedule func(*cron.Cron)) error {
	//加载配置文件
	conf, err := config.Load(confFile)
	if err != nil {
		return err
	}

	//初始化服务信息
	err = initServer()
	if err != nil {
		return fmt.Errorf("init server failed, %s", err.Error())
	}

	//容器初始化
	err = boot(conf)
	if err != nil {
		return fmt.Errorf("container ini failed %s", err.Error())
	}

	//注册Cron执行计划
	cronEngine := cron.New()
	registerSchedule(cronEngine)
	cronEngine.Start()

	//写pid文件
	writePidFile(pidFile)

	//注册信号量
	registerSignal()

	//等待停止信号
	waitConsoleStop(cronEngine)
	return nil
}
