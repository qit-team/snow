package main

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/app/http/routes"
	"github.com/qit-team/snow/app/console"
	"github.com/qit-team/snow/app/jobs"
	"github.com/qit-team/snow/bootstrap"
	"fmt"
	"os"
	"errors"
	"github.com/qit-team/snow-core/kernel/server"
	//启用本程序需要的各驱动
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/qit-team/snow-core/cache/rediscache"
	_ "github.com/qit-team/snow-core/queue/redisqueue"
	//_ "github.com/qit-team/snow-core/queue/alimnsqueue"
)

func main() {
	//解析启动命令
	opts := config.GetOptions()
	if opts.ShowVersion {
		fmt.Printf("%s\ncommit %s\nbuilt on %s\n", server.Version, server.BuildCommit, server.BuildDate)
		os.Exit(0)
	}

	handleCmd(opts)

	err := startServer(opts)
	if err != nil {
		fmt.Printf("server start error, %s\n", err)
		os.Exit(1)
	}
}

//执行(status|stop|restart)命令
func handleCmd(opts *config.Options) {
	if opts.Cmd != "" {
		pidFile := config.GenPidFile(opts)
		err := server.HandleUserCmd(opts.Cmd, pidFile)
		if err != nil {
			fmt.Printf("Handle user command(%s) error, %s\n", opts.Cmd, err)
		} else {
			fmt.Printf("Handle user command(%s) succ \n ", opts.Cmd)
		}
		os.Exit(0)
	}
}

func startServer(opts *config.Options) (err error) {
	//加载配置
	conf, err := config.Load(opts.ConfFile)
	if err != nil {
		return
	}

	//引导程序
	err = bootstrap.Bootstrap(conf)
	if err != nil {
		return
	}

	pidFile := config.GenPidFile(opts)

	//根据启动命令行参数，决定启动哪种服务模式
	switch opts.App {
	case "api":
		err = server.StartHttp(pidFile, conf.Api, routes.RegisterRoute)
	case "cron":
		err = server.StartConsole(pidFile, console.RegisterSchedule)
	case "job":
		err = server.StartJob(pidFile, jobs.RegisterWorker)
	default:
		err = errors.New("no server start")
	}
	return
}
