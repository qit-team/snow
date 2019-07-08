package alimnsqueue

import (
	"context"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/qit-team/snow/pkg/alimns"
	"strings"
	"errors"
)

const (
	DefaultVisibilityTimeout = int64(30)
)

type MnsQueueClient struct {
	client ali_mns.MNSClient
}

//实例模式
func NewMnsQueueClient(DiName string) (*MnsQueueClient) {
	m := new(MnsQueueClient)
	m.client = alimns.GetMns(DiName)
	return m
}

/**
 * 队列消息入队
 * args[0] delay 延迟消息，单位秒
 * args[1] priority
 */
func (m *MnsQueueClient) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	delay, priority := getOption(args...)

	//mns消息格式 可以设置优先级和延迟时间
	aliMsg := ali_mns.MessageSendRequest{
		MessageBody:  message,
		DelaySeconds: delay,
		Priority:     priority,
	}

	queueClient := alimns.GetMnsBasicQueue(m.client, key)
	_, err := queueClient.SendMessage(aliMsg)

	if err != nil {
		return false, err
	}

	return true, nil
}

/**
 * 队列消息出队
 * return 第一个参数是消息 第二个参数是mns的ReceiptHandle命名为token，通过token确定消息是否从队列删除
 */
func (m *MnsQueueClient) Dequeue(ctx context.Context, key string) (message string, token string, err error) {
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	//目前只做单次读取，不需要实现常驻进程，这部分由job完成

	//从alimns接收消息放入channel
	queueClient := alimns.GetMnsBasicQueue(m.client, key)

	go func() {
		queueClient.ReceiveMessage(respChan, errChan)
	}()

	select {
	case resp := <-respChan:
		//代表N秒内其他并发队列不可见这条消息
		if ret, err1 := queueClient.ChangeMessageVisibility(resp.ReceiptHandle, DefaultVisibilityTimeout); err1 != nil {
			err = err1
			return
		} else {
			//处理resp.MessageBody 阿里这什么sdk 也不说明各个函数作用。。。暂时就按照demo例子里用到的函数写了
			return resp.MessageBody, ret.ReceiptHandle, nil
		}
	case err2 := <-errChan:
		err = err2
		if strings.Contains(err2.Error(), "MessageNotExist") {
			//如果消息不存在的时候，返回的message为空字符串
			err = nil
			return
		}
	}
	return
}

/**
 * 队列消息批量入队
 * args[0] delay 延迟消息，单位秒
 * args[1] priority
 */
func (m *MnsQueueClient) BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) (bool, error) {
	if len(messages) == 0 {
		return false, errors.New("messages is empty")
	}

	delay, priority := getOption(args...)

	//mns消息格式 可以设置优先级和延迟时间
	msgArr := make([]ali_mns.MessageSendRequest, len(messages))
	for k, message := range messages {
		msgArr[k] = ali_mns.MessageSendRequest{
			MessageBody:  message,
			DelaySeconds: delay,
			Priority:     priority,
		}
	}

	queueClient := alimns.GetMnsBasicQueue(m.client, key)
	_, err := queueClient.BatchSendMessage(msgArr...)

	if err != nil {
		return false, err
	}

	return true, nil
}

/**
 * 确认消息接收
 */
func (m *MnsQueueClient) AckMsg(ctx context.Context, key string, token string) (bool, error) {
	queueClient := alimns.GetMnsBasicQueue(m.client, key)
	if len(token) < 1 {
		return false, errors.New("token empty")
	}
	err := queueClient.DeleteMessage(token)
	if err != nil {
		return false, err
	}
	return true, nil
}

//入队参数
func getOption(args ...interface{}) (delay int64, priority int64) {
	delay = 0
	priority = 1

	l := len(args)
	if l > 0 {
		de, ok := args[0].(int64)
		if ok {
			delay = de
		}

		if l > 1 {
			pr, ok := args[1].(int64)
			if ok {
				priority = pr
			}
		}
	}
	return
}
