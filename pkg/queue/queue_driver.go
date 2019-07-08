package queue

import "context"

//队列驱动接口，所有队列驱动都需要实现以下接口
type QueueDriver interface {
	//单入队
	Enqueue(ctx context.Context, key string, message string, args ...interface{}) (isOk bool, err error)
	//单出队： 消息不存在是返回空字符串
	Dequeue(ctx context.Context, key string) (message string, token string, err error)
	//确认接收消息redis用不到，alimns需要，后续可以接入kafka或者rabbitmq
	AckMsg(ctx context.Context, key string, token string) (isOk bool, err error)
	//单key批量入队
	BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) (isOk bool, err error)
}
