package main

const (
	_tplBuildBin = `*
!.gitignore
`

	_tplBuildShell = `#/bin/bash
os=$1 #系统linux
arch=$2 #架构amd64

#回到根目录
rootPath=$(cd ` + "`dirname $0`" + `/../../; pwd)

#编译
GOOS=$os GOARCH=$arch go build -o build/bin/snow main.go
`
)
