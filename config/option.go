package config

import (
	"flag"
	"strings"
)

var options *Options

//------------------------启动命令配置
type Options struct {
	ShowVersion bool
	Cmd         string
	ConfFile    string
	App         string
	PidDir      string
	Queue       string
	Command     string
}

func parseOptions() *Options {
	opts := new(Options)
	flag.BoolVar(&opts.ShowVersion, "v", false, "show version")
	flag.StringVar(&opts.App, "a", "api", "application to run")
	flag.StringVar(&opts.Cmd, "k", "", "status|stop|restart")
	flag.StringVar(&opts.ConfFile, "c", ".env", "conf file path")
	flag.StringVar(&opts.PidDir, "p", "/var/run/", "pid directory")
	flag.StringVar(&opts.Queue, "queue", "", "topic of queue is enable")
	flag.StringVar(&opts.Command, "m", "", "command name")
	flag.Parse()
	return opts
}

//获取启动命令配置
func GetOptions() *Options {
	if options == nil {
		options = parseOptions()
	}
	return options
}

//pid进程号的保存路径
func (opts *Options) GenPidFile() string {
	return strings.TrimRight(opts.PidDir, "/") + "/" + opts.App + ".pid"
}
