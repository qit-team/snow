package redisqueue

import (
	"context"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow/pkg/redis"
	"errors"
)

type RedisQueueClient struct {
	client *redis_pool.ReplicaPool
}

func NewRedisQueueClient(diName string) (*RedisQueueClient) {
	m := new(RedisQueueClient)
	m.client = redis.GetRedis(diName)
	return m
}

/**
 * 队列消息入队
 */
func (m *RedisQueueClient) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	//redis暂时不要延迟和优先级
	_, err := m.client.RPush(key, message)
	if err != nil {
		return false, err
	}
	return true, err
}

/**
 * 队列消息出队
 */
func (m *RedisQueueClient) Dequeue(ctx context.Context, key string) (message string, token string, err error) {
	message, err = m.client.LPop(key)
	if err == redis_pool.ErrNil {
		err = nil
		message = ""
	}
	return
}

/**
 * 确认消息接收 redis暂时用不到
 */
func (m *RedisQueueClient) AckMsg(ctx context.Context, key string, token string) (bool, error) {
	return true, nil
}

/**
 * 队列消息入队
 */
func (m *RedisQueueClient) BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) (bool, error) {
	//redis暂时不要延迟和优先级
	if len(messages) == 0 {
		return false, errors.New("messages is empty")
	}
	_, err := m.client.RPush(key, arrayStringToInterface(messages)...)
	if err != nil {
		return false, err
	}
	return true, err
}

func arrayStringToInterface(arr []string) []interface{} {
	newArr := make([]interface{}, len(arr))
	for k, v := range arr {
		newArr[k] = v
	}
	return newArr
}
