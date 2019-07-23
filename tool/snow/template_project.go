package main

const (
	_tplReadme = ``

	_tplGitignore = ``

	_tplGoMod = `module {{.ModuleName}}

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/qit-team/snow-core v0.1.5
	github.com/qit-team/work v0.3.3
	github.com/robfig/cron v1.2.0
)
`

	_tplMain = `package main

import (
	"{{.ModuleName}}/config"
	"{{.ModuleName}}/app/http/routes"
	"{{.ModuleName}}/app/console"
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/bootstrap"
	"{{.ModuleName}}/app/console/commands"
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
		pidFile := opts.GenPidFile()
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

	pidFile := opts.GenPidFile()

	//根据启动命令行参数，决定启动哪种服务模式
	switch opts.App {
	case "api":
		err = server.StartHttp(pidFile, conf.Api, routes.RegisterRoute)
	case "cron":
		err = server.StartConsole(pidFile, console.RegisterSchedule)
	case "job":
		err = server.StartJob(pidFile, jobs.RegisterWorker)
	case "command":
		err = server.ExecuteCommand(opts.Command, commands.RegisterCommand)
	default:
		err = errors.New("no server start")
	}
	return
}
`

	_tplEnv = `# toml配置文件
# Wiki：https://github.com/toml-lang/toml
Debug = true
Env = "local" # local-本地 develop-开发 beta-预发布 production-线上

[Log]
Handler = "file"
Dir = "./logs"
Level = "info"

[Db]
Driver = "mysql"

[Db.Option]
MaxConns = 128
MaxIdle = 32
IdleTimeout = 180 # second
Charset = "utf8mb4"
ConnectTimeout = 3 # second

[Db.Master]
Host = "127.0.0.1"
Port = 3306
User = "root"
Password = "123456"
DBName = "test"

[[Db.Slaves]] # 支持多个从库
Host = "127.0.0.1"
Port = 3306
User = "root"
Password = "123456"
DBName = "test"

[Api]
Host = "0.0.0.0"
Port = 8000

[Cache]
Driver = "redis"

[Redis.Master]
Host = "127.0.0.1"
Port = 6379
#Password = ""
#DB = 0

#[Redis.Option]
#MaxIdle = 64
#MaxConns = 256
#IdleTimeout = 180 # second
#ConnectTimeout = 1
#ReadTimeout = 1
#WriteTimeout = 1

[AliMns]
Url =  ""
AccessKeyId = ""
AccessKeySecret = ""
`

	_tplLog = `*
!.gitignore
`
)
