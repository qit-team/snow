package main

const (
	_tplReadme = `## Snow
Snow是一套简单易用的Go语言业务框架，整体逻辑设计简洁，支持HTTP服务、队列调度和任务调度等常用业务场景模式。

## Quick start

### Build
sh build/shell/build.sh

### Run
` + "```" + `shell
1. build/bin/snow -a api  #启动Api服务
2. build/bin/snow -a cron #启动Cron定时任务服务
3. build/bin/snow -a job  #启动队列调度服务
4. build/bin/snow -a command -m test  #执行名称为test的脚本任务
` + "```" + `

## Documents

- [项目地址](https://github.com/qit-team/snow)
- [中文文档](https://github.com/qit-team/snow/wiki)
- [changelog](https://github.com/qit-team/snow/blob/master/CHANGLOG.md)
- [xorm](http://gobook.io/read/github.com/go-xorm/manual-zh-CN/)
`

	_tplGitignore = `/.idea
/vendor
/.env
!/.env.example
`

	_tplGoMod = `module {{.ModuleName}}

go 1.12

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/SkyAPM/go2sky v0.6.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/fatih/color v1.13.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/go-resty/resty/v2 v2.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/qit-team/snow-core v0.1.28
	github.com/qit-team/work v0.3.11
	github.com/robfig/cron v1.2.0
	github.com/swaggo/gin-swagger v1.3.3
	github.com/swaggo/swag v1.7.6
	github.com/urfave/cli/v2 v2.3.0
	gopkg.in/go-playground/validator.v9 v9.31.0
	xorm.io/core v0.7.3
	xorm.io/xorm v1.2.5
)
`

	_tplMain = `package main

import (
	"errors"
	"fmt"
	"os"

	"{{.ModuleName}}/app/console"
	"{{.ModuleName}}/app/http/routes"
	"{{.ModuleName}}/app/jobs"
	"{{.ModuleName}}/bootstrap"
	"{{.ModuleName}}/config"
	_ "{{.ModuleName}}/docs"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/qit-team/snow-core/cache/rediscache"
	"github.com/qit-team/snow-core/kernel/server"
	_ "github.com/qit-team/snow-core/queue/redisqueue"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
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
		err = server.ExecuteCommand(opts.Command, console.RegisterCommand)
	default:
		err = errors.New("no server start")
	}
	return
}
`

	_tplEnv = `# toml配置文件
# Wiki：https://github.com/toml-lang/toml
ServiceName = "snow"
Debug = true
Env = "local" # local-本地 develop-开发 beta-预发布 production-线上
PrometheusCollectEnable = true
SkyWalkingOapServer = "127.0.0.1:11800"

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
Port = 8080

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
