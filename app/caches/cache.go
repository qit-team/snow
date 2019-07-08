package caches

import (
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/cache"
	"sync"
	"github.com/qit-team/snow/pkg/redis"
)

const (
	DiName = redis.SingletonMain //默认缓存依赖的实例别名
	Prefix = ""                  //默认缓存key前缀
	TTL    = 86400               //默认缓存时间
)

var cacheMu sync.Mutex

//缓存基类
type BaseCache struct {
	cache      *cache.Cache
	DiName     string //缓存依赖的实例别名
	Prefix     string //缓存key前缀
	DriverType string //缓存驱动
	TTL        int    //缓存时间
}

func (m *BaseCache) SetPrefix(prefix string) {
	m.Prefix = prefix
}

func (m *BaseCache) GetPrefix() string {
	return m.Prefix
}

func (m *BaseCache) GetPrefixOrDefault() string {
	if m.Prefix != "" {
		return m.Prefix
	} else {
		return Prefix
	}
}

func (m *BaseCache) SetDiName(diName string) {
	m.DiName = diName
}

func (m *BaseCache) GetDiName() string {
	return m.DiName
}

func (m *BaseCache) GetDiNameOrDefault() string {
	if m.DiName != "" {
		return m.DiName
	} else {
		return DiName
	}
}

func (m *BaseCache) SetDriverType(driverType string) {
	m.DriverType = driverType
}

func (m *BaseCache) GetDriverType() string {
	return m.DriverType
}

func (m *BaseCache) GetDriverTypeOrDefault() string {
	if m.DriverType != "" {
		return m.DriverType
	} else {
		return config.GetConf().Cache.Driver
	}
}

func (m *BaseCache) SeTTL(ttl int) {
	m.TTL = ttl
}

func (m *BaseCache) GetTTL() int {
	return m.TTL
}

func (m *BaseCache) GetTTLOrDefault() int {
	if m.TTL != 0 {
		return m.TTL
	} else {
		return TTL
	}
}

//获取缓存类
func (m *BaseCache) GetCache() *cache.Cache {
	//不使用once.Done是因为会有多种cache实例
	if m.cache == nil {
		cacheMu.Lock()
		defer cacheMu.Unlock()
		if m.cache == nil {
			diName := m.GetDiNameOrDefault()
			driverType := m.GetDriverTypeOrDefault()
			prefix := m.GetPrefixOrDefault()
			ttl := m.GetTTLOrDefault()
			m.cache, _ = cache.NewCache(diName, driverType, prefix, ttl)
		}
	}
	return m.cache
}
