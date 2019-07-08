package queue

import (
	"context"
	"github.com/qit-team/snow/pkg/queue/redisqueue"
	"github.com/qit-team/snow/pkg/queue/alimnsqueue"
	"sync"
	"errors"
)

const (
	DriverTypeRedis = "redis"
	DriverTypeMNS   = "ali_mns"
)

var (
	ErrDriverNil = errors.New("driver is nil")
)

type Queue struct {
	driver     QueueDriver
	driverType string //缓存驱动类型 目前支持redis和alimns
	queueName  string //队列名称
}

var (
	instances map[string]*Queue
	lock      sync.RWMutex
)

func init() {
	instances = make(map[string]*Queue)
}

//单例模式
func GetInstance(diName string, driverType string) (queue *Queue, err error) {
	key := diName + ":" + driverType
	lock.RLock()
	queue, ok := instances[key]
	lock.RUnlock()
	if ok {
		if queue != nil {
			return queue, nil
		}
	}

	queue, err = NewQueue(diName, driverType)
	if err != nil {
		return
	}

	if queue != nil {
		lock.Lock()
		instances[key] = queue
		lock.Unlock()
	}

	return
}

func NewQueue(diName string, driverType string) (*Queue, error) {
	q := new(Queue)
	//q.queueName = queueName
	switch driverType {
	case DriverTypeRedis:
		q.driverType = DriverTypeRedis
		//redis驱动
		q.driver = redisqueue.NewRedisQueueClient(diName)
		break
	case DriverTypeMNS:
		q.driverType = DriverTypeMNS
		//alimns驱动
		q.driver = alimnsqueue.NewMnsQueueClient(diName)
		break
	default:
	}

	if q.driver == nil {
		return nil, ErrDriverNil
	}

	return q, nil
}

func (q *Queue) SetDriver(driver QueueDriver) *Queue {
	q.driver = driver
	return q
}

func (q *Queue) GetDriver() QueueDriver {
	return q.driver
}

/*
 * 目前alimns驱动的话支持
 *  1、args[0] delay 延迟
 *  2、args[1] priority 优先级
 */
func (q *Queue) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	return q.driver.Enqueue(ctx, key, message, args...)
}

func (q *Queue) Dequeue(ctx context.Context, key string) (string, string, error) {
	return q.driver.Dequeue(ctx, key)
}

func (q *Queue) AckMsg(ctx context.Context, key string, token string) (bool, error) {
	return q.driver.AckMsg(ctx, key, token)
}

func (q *Queue) BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) (bool, error) {
	return q.driver.BatchEnqueue(ctx, key, messages, args...)
}
