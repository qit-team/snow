package jobs

import (
	"github.com/qit-team/work"
	"github.com/qit-team/snow/pkg/queue"
	"github.com/qit-team/snow/pkg/log/logger"
	"github.com/qit-team/snow/config"
	"strings"
	"context"
	"github.com/qit-team/snow/pkg/redis"
)

var (
	jb *work.Job
)

/**
 * 配置队列任务
 */
func RegisterWorker(job *work.Job) {
	setJob(job)

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
	q, err := queue.GetInstance(redis.SingletonMain, queue.DriverTypeRedis)
	if err != nil {
		panic("queue service init error:" + err.Error())
	}
	//针对topic设置相关的queue
	job.AddQueue(q, "topic-test1", "topic-test2")
	//设置默认的queue, 没有设置过的topic会使用默认的queue
	job.AddQueue(q)
}

/**
 * 设置配置参数
 */
func SetOptions(job *work.Job) {
	//设置logger，需要实现work.Logger接口的方法
	job.SetLogger(logger.GetLogger())

	//设置启用的topic，未设置表示启用全部注册过topic
	if config.GetOptions().Queue != "" {
		topics := strings.Split(config.GetOptions().Queue, ",")
		job.SetEnableTopics(topics...)
	}
}

func setJob(job *work.Job) {
	if jb == nil {
		jb = job
	}
}

func getJob() *work.Job {
	if jb == nil {
		jb = work.New()
		RegisterWorker(jb)
	}
	return jb
}

/**
 * 消息入队 -- 原始message
 */
func Enqueue(ctx context.Context, topic string, message string, args ...interface{}) (isOk bool, err error) {
	return getJob().Enqueue(ctx, topic, message, args...)
}

/**
 * 消息入队 -- Task数据结构
 */
func EnqueueWithTask(ctx context.Context, topic string, task work.Task, args ...interface{}) (isOk bool, err error) {
	return getJob().EnqueueWithTask(ctx, topic, task, args...)
}

/**
 * 消息批量入队 -- 原始message
 */
func BatchEnqueue(ctx context.Context, topic string, messages []string, args ...interface{}) (isOk bool, err error) {
	return getJob().BatchEnqueue(ctx, topic, messages, args...)
}

/**
 * 消息批量入队 -- Task数据结构
 */
func BatchEnqueueWithTask(ctx context.Context, topic string, tasks []work.Task, args ...interface{}) (isOk bool, err error) {
	return getJob().BatchEnqueueWithTask(ctx, topic, tasks, args...)
}
