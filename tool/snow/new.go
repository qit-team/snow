package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path"
)

func runNew(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return errors.New("required project name")
	}
	p.Name = ctx.Args().First()

	if p.ModuleName == "" {
		p.ModuleName = p.Name
	}

	if p.Path != "" {
		p.Path = path.Join(p.Path, p.Name)
	} else {
		pwd, _ := os.Getwd()
		p.Path = path.Join(pwd, p.Name)
	}
	// creata a project
	if err := create(); err != nil {
		return err
	}

	// 创建完项目执行安装swag命令
	if err := installSwag(); err != nil {
		fmt.Printf("Swagger install fail\n", err)
	}

	fmt.Printf("Project: %s\n", p.Name)
	fmt.Printf("Module Name: %s\n", p.ModuleName)
	fmt.Printf("Directory: %s\n\n", p.Path)
	fmt.Println("The application has been created.")
	return nil
}
