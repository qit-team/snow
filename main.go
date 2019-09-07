package main

import (
	"errors"
	"fmt"
	"github.com/qit-team/snow-core/kernel/server"
	"github.com/qit-team/snow/app/console"
	"github.com/qit-team/snow/app/http/routes"
	"github.com/qit-team/snow/app/jobs"
	"github.com/qit-team/snow/bootstrap"
	"github.com/qit-team/snow/config"
	"os"
	//启用本程序需要的各驱动
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/qit-team/snow-core/cache/rediscache"
	_ "github.com/qit-team/snow-core/queue/redisqueue"
	_ "github.com/qit-team/snow/docs"
	//_ "github.com/qit-team/snow-core/queue/alimnsqueue"
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
