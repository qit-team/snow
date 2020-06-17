package main

const (
	_tplJobKernel = `package jobs

import (
	"strings"

	"{{.ModuleName}}/app/jobs/basejob"

	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow-core/queue"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/work"
)

var enableQueues string

func SetEnableQueue(q string) {
	enableQueues = q
}

/**
 * 配置队列任务
 */
func RegisterWorker(job *work.Job) {
	basejob.SetJob(job)

	//设置worker的任务投递回调函数
	job.AddFunc("topic-test", test)
	//设置worker的任务投递回调函数，和并发数
	job.AddFunc("topic-test1", test, 2)
	//使用worker结构进行注册
	job.AddWorker("topic-test2", &work.Worker{Call: work.MyWorkerFunc(test), MaxConcurrency: 1})

	RegisterQueueDriver(job)
	SetOptions(job)
}

/**
 * 给topic注册对应的队列服务
 */
func RegisterQueueDriver(job *work.Job) {
	//设置队列服务，需要实现work.Queue接口的方法
	q := queue.GetQueue(redis.SingletonMain, queue.DriverTypeRedis)
	//针对topic设置相关的queue
	job.AddQueue(q, "topic-test1", "topic-test2")
	//设置默认的queue, 没有设置过的topic会使用默认的queue
	job.AddQueue(q)
}

/**
 * 设置配置参数
 */
func SetOptions(job *work.Job) {
	// 设置logger，需要实现work.Logger接口的方法
	job.SetLogger(logger.GetLogger())

	// 设置启用的topic，未设置表示启用全部注册过topic
	if enableQueues != "" {
		topics := strings.Split(enableQueues, ",")
		job.SetEnableTopics(topics...)
	}
}
`

	_tplJobTest = `package jobs

import (
	"fmt"
	"time"

	"github.com/qit-team/work"
)

func test(task work.Task) (work.TaskResult) {
	time.Sleep(time.Millisecond * 5)
	s, err := work.JsonEncode(task)
	if err != nil {
		//work.StateFailed 不会进行ack确认
		//work.StateFailedWithAck 会进行actk确认
		//return work.TaskResult{Id: task.Id, State: work.StateFailed}
		return work.TaskResult{Id: task.Id, State: work.StateFailedWithAck}
	} else {
        //work.StateSucceed 会进行ack确认
		fmt.Println("do task", s)
		return work.TaskResult{Id: task.Id, State: work.StateSucceed}
	}

}
`

	_tplJobBase = `package basejob

import (
	"context"
	"sync"

	"github.com/qit-team/work"
)

var (
	jb       *work.Job
	register func(job *work.Job)
	mu       sync.RWMutex
)

func SetJob(job *work.Job) {
	if jb == nil {
		jb = job
	}
}

func SetJobRegister(r func(*work.Job)) {
	register = r
}

func GetJob() *work.Job {
	if jb == nil {
		if register != nil {
			mu.Lock()
			defer mu.Unlock()
			jb = work.New()
			register(jb)
		} else {
			panic("job register is nil")
		}
	}
	return jb
}

/**
 * 消息入队 -- 原始message
 */
func Enqueue(ctx context.Context, topic string, message string, args ...interface{}) (isOk bool, err error) {
	return GetJob().Enqueue(ctx, topic, message, args...)
}

/**
 * 消息入队 -- Task数据结构
 */
func EnqueueWithTask(ctx context.Context, topic string, task work.Task, args ...interface{}) (isOk bool, err error) {
	return GetJob().EnqueueWithTask(ctx, topic, task, args...)
}

/**
 * 消息批量入队 -- 原始message
 */
func BatchEnqueue(ctx context.Context, topic string, messages []string, args ...interface{}) (isOk bool, err error) {
	return GetJob().BatchEnqueue(ctx, topic, messages, args...)
}

/**
 * 消息批量入队 -- Task数据结构
 */
func BatchEnqueueWithTask(ctx context.Context, topic string, tasks []work.Task, args ...interface{}) (isOk bool, err error) {
	return GetJob().BatchEnqueueWithTask(ctx, topic, tasks, args...)
}
`
)
