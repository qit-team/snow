package config

import "flag"

var options *Options

//------------------------启动命令配置
type Options struct {
	ShowVersion bool
	Cmd         string
	ConfFile    string
	App         string
	PidPath     string
	Queue       string
}

func parseOptions() *Options {
	opts := new(Options)
	flag.BoolVar(&opts.ShowVersion, "v", false, "Show Version")
	flag.StringVar(&opts.App, "a", "api", "application to run")
	flag.StringVar(&opts.Cmd, "k", "", "status|stop|restart")
	flag.StringVar(&opts.ConfFile, "c", ".env", "conf file path")
	flag.StringVar(&opts.PidPath, "p", "/var/run/", "pid file path")
	flag.StringVar(&opts.Queue, "queue", "", "topic of queue is enable")
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
