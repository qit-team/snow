package rediscache

import (
	"context"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow/pkg/redis"
)

type RedisCacheClient struct {
	client *redis_pool.ReplicaPool
}

//实例模式
func NewRedisCacheClient(diName string) (*RedisCacheClient) {
	m := new(RedisCacheClient)
	m.client = redis.GetRedis(diName)
	return m
}

/**
 * 获取缓存key的数据
 * 注意事项，如果key值不存在的话，返回的是空字符串，而不是nil
 */
func (m *RedisCacheClient) Get(ctx context.Context, key string) (interface{}, error) {
	return m.client.Get(key)
}

func (m *RedisCacheClient) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	cKeys := m.convert(keys)
	values, err := m.client.MGet(cKeys...)
	if err != nil {
		return nil, err
	}

	arr := make(map[string]interface{})

	for index, key := range keys {
		arr[key] = values[index]
	}
	return arr, nil
}

func (m *RedisCacheClient) Set(ctx context.Context, key string, value interface{}, ttl int) (bool, error) {
	return m.client.SetEX(key, value, int64(ttl))
}

func (m *RedisCacheClient) SetMulti(ctx context.Context, items map[string]interface{}, ttl int) (bool, error) {
	return m.client.MSet(items)
}

func (m *RedisCacheClient) Delete(ctx context.Context, key string) (bool, error) {
	res, err := m.client.Del(key)
	return res > 0, err
}

func (m *RedisCacheClient) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	cKeys := m.convert(keys)
	res, err := m.client.Del(cKeys...)
	return res > 0, err
}

func (m *RedisCacheClient) Expire(ctx context.Context, key string, ttl int) (bool, error) {
	return m.client.Expire(key, int64(ttl))
}

func (m *RedisCacheClient) convert(keys []string) []interface{} {
	t := make([]interface{}, len(keys))
	for i, v := range keys {
		t[i] = v
	}
	return t
}
