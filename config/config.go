package config

import (
	"github.com/BurntSushi/toml"
	"github.com/qit-team/snow-core/config"
	"os"
)

const (
	ProdEnv  = "production" // 线上环境
	BetaEnv  = "beta"       // beta环境
	DevEnv   = "develop"    // 开发环境
	LocalEnv = "local"      // 本地环境
)

var srvConf *Config

// ------------------------配置文件解析
type Config struct {
	Env   string             `toml:"Env"`
	Debug bool               `toml:"Debug"`
	Log   config.LogConfig   `toml:"Log"`
	Redis config.RedisConfig `toml:"Redis"`
	Mns   config.MnsConfig   `toml:"AliMns"`
	Db    config.DbConfig    `toml:"Db"`
	Api   config.ApiConfig   `toml:"Api"`
}

// ------------------------ 加载配置 ------------------------//
func Load(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	conf := new(Config)
	if _, err := toml.DecodeFile(path, conf); err != nil {
		return nil, err
	}
	srvConf = conf
	return conf, nil
}

// 当前配置
func GetConf() *Config {
	return srvConf
}

// 是否调试模式
func IsDebug() bool {
	return srvConf.Debug
}

// 当前环境，默认本地开发
func GetEnv() string {
	if srvConf.Env == "" {
		return LocalEnv
	}
	return srvConf.Env
}

// 是否当前环境
func IsEnvEqual(env string) bool {
	return GetEnv() == env
}
