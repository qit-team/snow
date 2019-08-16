package bannerlistcache

import (
	"github.com/qit-team/snow/app/caches"
	"sync"
	"github.com/qit-team/snow-core/cache"
)

const (
	prefix = caches.BannerList // 缓存key的前缀
)

var (
	instance *bannerListCache
	once     sync.Once
)

type bannerListCache struct {
	cache.BaseCache
}

// 单例模式
func GetInstance() *bannerListCache {
	once.Do(func() {
		instance = new(bannerListCache)
		instance.Prefix = prefix
		// instance.DiName = redis.SingletonMain // 设置缓存依赖的实例别名
		// instance.DriverType = cache.DriverTypeRedis // 设置缓存驱动的类型,默认redis
		// instance.SeTTL(86400) 设置默认缓存时间 默认86400
	})
	return instance
}
