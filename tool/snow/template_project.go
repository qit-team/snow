package main

const (
	_tplReadme = `## Snow
Snow是一套简单易用的Go语言业务框架，整体逻辑设计简洁，支持HTTP服务、队列调度和任务调度等常用业务场景模式。

## Quick start

### Build
sh build/shell/build.sh

### Run
` + "```" + `shell
1. build/bin/{{.ModuleName}} a api  #启动Api服务
2. build/bin/{{.ModuleName}} a cron #启动Cron定时任务服务
3. build/bin/{{.ModuleName}} a job -queue test  #启动队列调度服务
4. build/bin/{{.ModuleName}} a command -m test  #执行名称为test的脚本任务
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
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/ouqiang/goutil v1.2.3
	github.com/prometheus/client_golang v1.6.0
	github.com/qit-team/snow-core v0.1.19
	github.com/qit-team/work v0.3.9
	github.com/robfig/cron v1.2.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/urfave/cli v1.22.4
	github.com/valyala/fasttemplate v1.1.0 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	golang.org/x/tools v0.0.0-20200601175630-2caf76543d99 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
	xorm.io/core v0.7.3
	xorm.io/xorm v1.0.2
)
`

	_tplMain = `package main

import (
	"log"
	"os"

	"{{.ModuleName}}/cli"
	_ "{{.ModuleName}}/docs"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/qit-team/snow-core/cache/rediscache"
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
	if err := cli.GetApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
