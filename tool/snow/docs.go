package main

import (
	"errors"
	"os"
	"github.com/urfave/cli"
	"fmt"
)

var (
	d *docs
)

type docs struct {
	Path       string //项目目录
}


func init() {
	d = new(docs)
}

// generate swag doc
func runDoc(ctx *cli.Context) (err error) {

	if d.Path == "" {
		d.Path, _ = os.Getwd()
	} else {
		if !isDirExist(d.Path) {
			return errors.New("project directory is not exist")
		}
	}

	err = runSwagInit(d.Path)
	if err == nil {
		fmt.Println("snow generate doc success!")
	} else {
		fmt.Println("snow generate doc error:", err)
	}
	return err
}