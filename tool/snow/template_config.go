package main

const (
	_tplConfig = `package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"github.com/qit-team/snow-core/config"
)

const (
	ProdEnv  = "production" //线上环境
	BetaEnv  = "beta"       //beta环境
	DevEnv   = "develop"    //开发环境
	LocalEnv = "local"      //本地环境
)

var srvConf *Config

//------------------------配置文件解析
type Config struct {
	Env   string             ` + "`toml:\"Env\"`" + `
	Debug bool               ` + "`toml:\"Debug\"`" + `
	Log   config.LogConfig   ` + "`toml:\"Log\"`" + `
	Redis config.RedisConfig ` + "`toml:\"Redis\"`" + `
	Mns   config.MnsConfig   ` + "`toml:\"AliMns\"`" + `
	Db    config.DbConfig    ` + "`toml:\"Db\"`" + `
	Api   config.ApiConfig   ` + "`toml:\"Api\"`" + `
}

func newConfig() *Config {
	return new(Config)
}

//------------------------ 加载配置 ------------------------//
func Load(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	conf := newConfig()
	if _, err := toml.DecodeFile(path, conf); err != nil {
		return nil, err
	}
	srvConf = conf
	return conf, nil
}

//当前配置
func GetConf() *Config {
	return srvConf
}

//是否调试模式
func IsDebug() bool {
	return srvConf.Debug
}

//当前环境，默认本地开发
func GetEnv() string {
	if srvConf.Env == "" {
		return LocalEnv
	}
	return srvConf.Env
}

//是否当前环境
func IsEnvEqual(env string) bool {
	return GetEnv() == env
}
`

	_tplOption = `package config

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
`
)
