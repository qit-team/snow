package redis

import (
	"errors"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow/config"
	"time"
)

//redis连接池实例，不对外暴露，通过redis_service_provider实现依赖注入和资源获取
func NewRedisClient(redisConf config.RedisConfig) (*redis_pool.ReplicaPool, error) {
	if redisConf.Master.Host == "" {
		return nil, errors.New("redis config is empty")
	}

	replicaConfig := &redis_pool.ReplicaConfig{
		Master: genRedisConfig(redisConf.Master),
		Slaves: []redis_pool.RedisConfig{},
		Opts:   genOptions(redisConf.Option),
	}

	for _, v := range redisConf.Slaves {
		replicaConfig.Slaves = append(replicaConfig.Slaves, genRedisConfig(v))
	}

	pool := redis_pool.NewReplicaPool(replicaConfig)
	return pool, nil
}

func genRedisConfig(c config.RedisBaseConfig) redis_pool.RedisConfig {
	return redis_pool.RedisConfig{
		Host:     c.Host,
		Port:     c.Port,
		Password: c.Password,
		DB:       c.DB,
	}
}

func genOptions(c config.RedisOptionConfig) redis_pool.Options {
	return redis_pool.Options{
		MaxIdle:        c.MaxIdle,
		MaxActive:      c.MaxConns,
		Wait:           c.Wait,
		IdleTimeout:    c.IdleTimeout * time.Second,
		ConnectTimeout: c.ConnectTimeout * time.Second,
		ReadTimeout:    c.ReadTimeout * time.Second,
		WriteTimeout:   c.WriteTimeout * time.Second,
	}
}
