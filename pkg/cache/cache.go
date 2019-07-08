package cache

import (
	"context"
	"github.com/qit-team/snow/app/utils"
	"github.com/qit-team/snow/pkg/cache/rediscache"
	"errors"
)

const (
	DriverTypeRedis = "redis"
)

var (
	ErrDriverNil = errors.New("driver is nil")
)

type Cache struct {
	driver     Driver
	driverType string //缓存驱动类型 目前支持redis
	prefix     string //前缀key
	TTL        int    //缓存时间
}

// args columns: TTL int
func NewCache(diName string, driverType string, prefix string, args ...int) (*Cache, error) {
	c := new(Cache)
	c.prefix = prefix
	if len(args) > 0 {
		c.TTL = args[0]
	}

	switch driverType {
	case DriverTypeRedis:
		c.driverType = DriverTypeRedis
		c.driver = rediscache.NewRedisCacheClient(diName)
		break
	default:
	}

	if c.driver == nil {
		return nil, ErrDriverNil
	}

	return c, nil
}

func (c *Cache) SetDriver(driver Driver) *Cache {
	c.driver = driver
	return c
}

func (c *Cache) GetDriver() Driver {
	return c.driver
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	key = c.key(key)
	return c.driver.Get(ctx, key)
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl int) (bool, error) {
	key = c.key(key)
	return c.driver.Set(ctx, key, value, ttl)
}

func (c *Cache) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	keys = c.keys(keys...)
	items, err := c.driver.GetMulti(ctx, keys...)
	if err != nil {
		return nil, err
	}

	m2 := make(map[string]interface{})
	for key, val := range items {
		m2[c.removePrefix(key)] = val
	}
	return m2, nil
}

func (c *Cache) SetMulti(ctx context.Context, items map[string]interface{}, ttl int) (bool, error) {
	arr := make(map[string]interface{})
	for key, value := range items {
		key = c.key(key)
		arr[key] = value
	}
	return c.driver.SetMulti(ctx, arr, ttl)
}

func (c *Cache) Delete(ctx context.Context, key string) (bool, error) {
	key = c.key(key)
	return c.driver.Delete(ctx, key)
}

func (c *Cache) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	keys = c.keys(keys...)
	return c.driver.DeleteMulti(ctx, keys...)
}

func (c *Cache) Expire(ctx context.Context, key string, ttl int) (bool, error) {
	key = c.key(key)
	return c.driver.Expire(ctx, key, ttl)
}

//补全key
func (c *Cache) key(key string) string {
	return c.prefix + key
}

//批量补全
func (c *Cache) keys(keys ...string) []string {
	arr := make([]string, len(keys))
	for i, key := range keys {
		arr[i] = c.key(key)
	}
	return arr
}

//去除前缀
func (c *Cache) removePrefix(key string) string {
	l := len(c.prefix)
	keyLen := len(key)
	return utils.Substr(key, l, keyLen-l)
}
