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
					Usage:       "directory for create project, default: current position",
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
					Usage:       "project directory, default: current position",
					Destination: &m.Path,
				},
				cli.StringFlag{
					Name:        "table, t",
					Value:       "",
					Usage:       "table name for new model, default: model name",
					Destination: &m.Table,
				},
				cli.StringFlag{
					Name:        "dsn, d",
					Value:       "",
					Usage:       "database dsn config, default: GetEnv('SNOW_DSN') or 'root:123456@tcp(localhost:3306)/test'",
					Destination: &m.DSN,
				},
				cli.StringFlag{
					Name:        "db, b",
					Value:       "",
					Usage:       "database name, will replace dsn's database",
					Destination: &m.DB,
				},
				cli.StringFlag{
					Name:        "driver, r",
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
			Usage:  "snow self-upgrade",
			Action: upgradeAction,
		},
		{
			Name:    "doc",
			Aliases: []string{"d"},
			Usage:   "generate doc",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "path, p",
					Value:       "",
					Usage:       "project directory, default: current position",
					Destination: &d.Path,
				},
			},

			Action: runDoc,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
