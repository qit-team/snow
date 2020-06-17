package cli

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/qit-team/snow/app/console"
	"github.com/qit-team/snow/app/http/routes"
	"github.com/qit-team/snow/app/jobs"
	"github.com/qit-team/snow/bootstrap"
	"github.com/qit-team/snow/config"

	"github.com/ouqiang/goutil"
	"github.com/qit-team/snow-core/kernel/server"
	"github.com/urfave/cli"
)

var (
	AppVersion           = "1.0.0"
	BuildDate, GitCommit string
	app                  *cli.App
)

func init() {
	app = cli.NewApp()
	app.Usage = "snow service"
	app.Version, _ = goutil.FormatAppVersion(AppVersion, GitCommit, BuildDate)
	app.Commands = commands()
	app.Flags = flags()
	app.Before = before
}

// get Commands
func commands() []cli.Command {
	appCommand := cli.Command{
		Name:  "a",
		Usage: "application to run",
		Subcommands: []cli.Command{
			// api
			{
				Name:  "api",
				Usage: "run api server",
				Action: func(ctx *cli.Context) {
					pDir := ctx.GlobalString("p")
					pidFile := genPidFile("api", pDir)
					if config.IsDebug() {
						server.SetDebug(true)
					}
					if err := server.StartHttp(pidFile, config.GetConf().Api, routes.RegisterRoute); err != nil {
						log.Fatal(err)
					}
				},
			},
			// cron
			{
				Name:  "cron",
				Usage: "run cron server",
				Action: func(ctx *cli.Context) {
					pDir := ctx.GlobalString("p")
					pidFile := genPidFile("cron", pDir)
					if err := server.StartConsole(pidFile, console.RegisterSchedule); err != nil {
						log.Fatal(err)
					}
				},
			},
			// job
			{
				Name:  "job",
				Usage: "run job server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "queue",
						Usage: "topics of queue is enable",
					},
				},
				Action: func(ctx *cli.Context) {
					jobs.SetEnableQueue(ctx.String("queue"))
					pDir := ctx.GlobalString("p")
					pidFile := genPidFile("job", pDir)
					if err := server.StartJob(pidFile, jobs.RegisterWorker); err != nil {
						log.Fatal(err)
					}
				},
			},
			// command
			{
				Name:  "command",
				Usage: "run command server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "m",
						Usage: "command name",
					},
				},
				Action: func(ctx *cli.Context) {
					command := ctx.String("m")
					if err := server.ExecuteCommand(command, console.RegisterCommand); err != nil {
						log.Fatal(err)
					}
				},
			},
		},
	}
	cmdCommand := cli.Command{
		Name:  "k",
		Usage: "status|stop|restart",
		Action: func(ctx *cli.Context) {
			if ctx.NArg() == 0 {
				log.Fatalf("do not specified parameter 'status|stop|restart'")
			}
			cmd := ctx.Args()[0]
			pDir := ctx.GlobalString("p")
			pidFile := genPidFile("api", pDir)
			err := server.HandleUserCmd(cmd, pidFile)
			if err != nil {
				log.Fatalf("Handle user command(%s) error, %s", cmd, err)
			} else {
				log.Printf("Handle user command(%s) succ", cmd)
			}
			os.Exit(0)
		},
	}
	return []cli.Command{appCommand, cmdCommand}
}

// get Flags
func flags() []cli.Flag {
	confFlag := cli.StringFlag{
		Name:  "c",
		Usage: "conf file path",
		Value: ".env",
	}
	pidFlag := cli.StringFlag{
		Name:  "p",
		Usage: "pid directory",
		Value: "/var/run/",
	}
	return []cli.Flag{confFlag, pidFlag}
}

func before(ctx *cli.Context) error {
	confFile := ctx.GlobalString("c")

	//加载配置
	conf, err := config.Load(confFile)
	if err != nil {
		return err
	}

	//引导程序
	err = bootstrap.Bootstrap(conf)
	if err != nil {
		return err
	}

	return nil
}

func genPidFile(app, pidDir string) string {
	p := strings.TrimRight(pidDir, "/") + "/" + app + ".pid"
	_, err := os.Stat(p)
	if err == os.ErrNotExist {
		os.MkdirAll(path.Dir(p), 0755)
	}
	return p
}

func GetApp() *cli.App {
	return app
}
