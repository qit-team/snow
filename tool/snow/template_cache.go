package main

const (
	_tplCacheKey = `package caches

//缓存前缀key，不同的业务使用不同的前缀，避免了业务之间的重用冲突
const (
	Cookie     = "ck:"
	Copy       = "cp:"
	BannerList = "bl:"
)
`

	_tplBannerListCache = `package bannerlistcache

import (
	"sync"

	"{{.ModuleName}}/app/caches"

	"github.com/qit-team/snow-core/cache"
)

const (
	prefix = caches.BannerList //缓存key的前缀
)

var (
	instance *bannerListCache
	once     sync.Once
)

type bannerListCache struct {
	cache.BaseCache
}

//单例模式
func GetInstance() *bannerListCache {
	once.Do(func() {
		instance = new(bannerListCache)
		instance.Prefix = prefix
		//instance.DiName = redis.SingletonMain //设置缓存依赖的实例别名
		//instance.DriverType = cache.DriverTypeRedis //设置缓存驱动的类型,默认redis
		//instance.SeTTL(86400) 设置默认缓存时间 默认86400
	})
	return instance
}
`

	_tplBannerListCacheTest = `package bannerlistcache

import (
	"context"
	"fmt"
	"testing"

	"{{.ModuleName}}/config"

	"github.com/qit-team/snow-core/cache"
	_ "github.com/qit-team/snow-core/cache/rediscache"
	"github.com/qit-team/snow-core/redis"
)

func init() {
	//加载配置文件
	conf, err := config.Load("../../../.env")
	if err != nil {
		fmt.Println(err)
	}

	//注册redis类
	err = redis.Pr.Register(cache.DefaultDiName, conf.Redis)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_GetMulti(t *testing.T) {
	ctx := context.TODO()
	cache := GetInstance()
	res, _ := cache.Set(ctx, "1000", "a")
	if res != true {
		t.Errorf("set key:%s is error", "1000")
	}

	keys := []string{"1000", "-2000", "9999"}
	cacheList, err := cache.GetMulti(ctx, keys...)
	if err != nil {
		t.Errorf("getMulti error:%s", err.Error())
	}
	fmt.Println(cacheList)
	i := 0
	for k, v := range cacheList {
		i++
		if k == "1000" {
			if v != "a" {
				t.Errorf("value of key:%s is error %v", k, v)
			}
		} else {
			if v != "" {
				t.Errorf("value of key:%s is error %v", k, v)
			}
		}
	}
	if i != len(keys) {
		t.Errorf("count of cache key is error: %d", i)
	}
}
`
)
