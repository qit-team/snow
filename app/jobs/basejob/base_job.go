package basejob

import (
	"github.com/qit-team/work"
	"context"
	"sync"
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
