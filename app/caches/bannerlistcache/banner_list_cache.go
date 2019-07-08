package bannerlistcache

import (
	"github.com/qit-team/snow/app/caches"
	"github.com/qit-team/snow/pkg/cache"
	"sync"
)

const (
	prefix = caches.BannerList //缓存key的前缀
)

type bannerListCache struct {
	caches.BaseCache
}

var (
	instance *bannerListCache
	once     sync.Once
)

//单例模式
func GetInstance() *bannerListCache {
	once.Do(func() {
		instance = new(bannerListCache)
		instance.SetPrefix(prefix)
		//instance.SetDiName(?) 设置缓存依赖的实例别名 默认redis.SingletonMain
		//instance.SetDriverType(?) 设置缓存驱动的类型 默认驱动类型为env配置文件
		//instance.SeTTL(?) 设置缓存时间 默认86400
	})
	return instance
}

//获取缓存类
func GetCache() *cache.Cache {
	return GetInstance().GetCache()
}
