package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

const (
	ProdEnv  = "production" //线上环境
	BetaEnv  = "beta"       //beta环境
	DevEnv   = "develop"    //开发环境
	LocalEnv = "local"      //本地环境
)

var srvConf *Config

//------------------------配置文件解析

type CacheConfig struct {
	Driver string //驱动类型，目前支持redis
}

type RedisBaseConfig struct {
	Host     string
	Port     int
	Password string
	DB       int //第几个库，默认0
}

type RedisOptionConfig struct {
	MaxIdle        int
	MaxConns       int
	Wait           bool
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type RedisConfig struct {
	Master RedisBaseConfig
	Slaves []RedisBaseConfig
	Option RedisOptionConfig
}

type DbBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DbOptionConfig struct {
	MaxIdle        int
	MaxConns       int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	Charset        string
}

type DbConfig struct {
	Driver string //驱动类型，目前支持mysql、postgres、mssql、sqlite3
	Master DbBaseConfig
	Slaves []DbBaseConfig
	Option DbOptionConfig
}

type MnsConfig struct {
	Url             string
	AccessKeyId     string
	AccessKeySecret string
}

type Config struct {
	Debug      bool        `toml:"Debug"`
	LogHandler string      `toml:"LogHandler"`
	LogDir     string      `toml:"LogDir"`
	LogLevel   string      `toml:"LogLevel"`
	Env        string      `toml:"Env"`
	Cache      CacheConfig `toml:"cache"`
	Redis      RedisConfig `toml:"Redis"`
	Mns        MnsConfig   `toml:"AliMns"`
	Db         DbConfig    `toml:"Db"`

	Api struct {
		Host string
		Port int
	} `toml:"Api"`
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
