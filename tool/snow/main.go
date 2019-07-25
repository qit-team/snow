package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "snow"
	app.Usage = "snow工具集"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "create new project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       "",
					Usage:       "project directory for create project",
					Destination: &p.Path,
				},
				cli.StringFlag{
					Name:        "module, m",
					Usage:       "project module name, for go mod init",
					Destination: &p.ModuleName,
				},
			},
			Action: runNew,
		},
		{
			Name:    "model",
			Aliases: []string{"m"},
			Usage:   "create new model",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       "",
					Usage:       "project directory for new model, default: current directory",
					Destination: &m.Path,
				},
				cli.StringFlag{
					Name:        "table, t",
					Value:       "",
					Usage:       "table name for new model, default: model name",
					Destination: &m.Table,
				},
				cli.StringFlag{
					Name:        "dsn",
					Value:       "",
					Usage:       "database dsn config, eg. root:123345@localhost:3306/test?charset=utf8mb4",
					Destination: &m.DSN,
				},
				cli.StringFlag{
					Name:        "database, b",
					Value:       "",
					Usage:       "database name. when dsn is empty, it is valid.",
					Destination: &m.DB,
				},
				cli.StringFlag{
					Name:        "driver, d",
					Value:       "",
					Usage:       "driver type, default: mysql",
					Destination: &m.DB,
				},
			},
			Action: runModel,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "snow version",
			Action: func(c *cli.Context) error {
				fmt.Println(getVersion())
				return nil
			},
		},
		{
			Name:   "upgrade",
			Usage:  "snows self-upgrade",
			Action: upgradeAction,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
