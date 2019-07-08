package queue

import (
	"testing"
	"context"
	"github.com/qit-team/snow/config"
	"fmt"
	"github.com/qit-team/snow/pkg/alimns"
	"github.com/qit-team/snow/pkg/redis"
)

func init() {
	//加载配置文件
	conf, err := config.Load("../../.env")
	if err != nil {
		fmt.Println(err)
	}

	//注册redis类
	err = (&redis.RedisServiceProvider{}).Register(conf)
	if err != nil {
		fmt.Println(err)
	}

	err = (&alimns.AliMnsServiceProvider{}).Register(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_Set(t *testing.T) {

	//redis_enqueue("hello")
	//redis_enqueue("thank u")
	//
	//redis_dequeue()
	//redis_dequeue()
	//
	//mns_enqueue_with_delay("msg_delay", 3)
	//mns_enqueue("msg1")
	//mns_enqueue("msg2")
	//
	//mns_dequeue()
	//mns_dequeue()
	//mns_dequeue()
}

func redis_enqueue(msg string) {
	diName := redis.SingletonMain
	queueName := "snow-test-queue"
	ctx := context.TODO()
	client, _ := NewQueue(diName, "redis")

	flag, _ := client.Enqueue(ctx, queueName, msg)
	//t.Log(flag, reflect.TypeOf(flag))
	fmt.Println("===RedisEnqueueFlag", msg, flag)
}

func redis_dequeue() {
	diName := redis.SingletonMain
	queueName := "snow-test-queue"
	ctx := context.TODO()
	client, _ := NewQueue(diName, "redis")

	ret, _, _ := client.Dequeue(ctx, queueName)
	//t.Log(ret, reflect.TypeOf(ret))
	fmt.Println("===RedisDequeueRet", ret)
}

func mns_enqueue(msg string) {
	queueName := "snow-test-queue"
	diNameMns := alimns.SingletonMain
	ctx := context.TODO()
	mnsClient, _ := NewQueue(diNameMns, "ali_mns")

	flag, _ := mnsClient.Enqueue(ctx, queueName, msg)
	fmt.Println("===MnsEnqueueFlag", msg, flag)

}

func mns_dequeue() {
	queueName := "snow-test-queue"
	diNameMns := alimns.SingletonMain
	ctx := context.TODO()
	mnsClient, _ := NewQueue(diNameMns, "ali_mns")

	msg, token, _ := mnsClient.Dequeue(ctx, queueName)
	fmt.Println("===MnsDequeueRet", msg, token)

	//mns读取完消息 直接ack
	flagAck, _ := mnsClient.AckMsg(ctx, queueName, token)
	fmt.Println("===MnsDequeueAck", flagAck)
}

func mns_enqueue_with_delay(msg string, delay int64) {
	queueName := "snow-test-queue"
	diNameMns := alimns.SingletonMain
	ctx := context.TODO()
	mnsClient, _ := NewQueue(diNameMns, "ali_mns")

	flag, _ := mnsClient.Enqueue(ctx, queueName, msg, delay)
	fmt.Println("===MnsEnqueueFlag", msg, flag)

}

func TestQueue_Redis_Batch(t *testing.T) {
	DiName := redis.SingletonMain
	client, err := NewQueue(DiName, "redis")
	if err != nil {
		t.Error("new queue error", err)
	}

	key1 := "snow-test-queue1"
	messages := []string{"msg1", "msg2"}
	batchEnqueue(t, client, key1, messages)

	key2 := "snow-test-queue2"
	messages = []string{"msg3", "msg4"}
	batchEnqueue(t, client, key2, messages)

	dequeue(t, client, key1)
	dequeue(t, client, key1)
	dequeue(t, client, key2)
	dequeue(t, client, key2)
}

func TestQueue_Mns_Batch(t *testing.T) {
	return
	DiName := alimns.SingletonMain
	client, err := NewQueue(DiName, "ali_mns")
	if err != nil {
		t.Error("new queue error", err)
	}

	key1 := "snow-test-queue1"
	messages := []string{"msg1", "msg2"}
	batchEnqueue(t, client, key1, messages)

	key2 := "snow-test-queue2"
	messages = []string{"msg3", "msg4"}
	//batchEnqueue(t, client, key2, messages)

	dequeue(t, client, key1)
	dequeue(t, client, key1)
	dequeue(t, client, key2)
	dequeue(t, client, key2)
}

func batchEnqueue(t *testing.T, client *Queue, key string, messages []string) {
	ctx := context.TODO()
	_, err := client.BacthEnqueue(ctx, key, messages)
	if err != nil {
		t.Error("batch enqueue error", err)
	}
	fmt.Println("batch enqueue", key, messages)
}

func dequeue(t *testing.T, client *Queue, key string) {
	ctx := context.TODO()

	message, token, err := client.Dequeue(ctx, key)
	if err != nil {
		t.Error("batch dequeue error")
	}

	if token != "" {
		_, err = client.AckMsg(ctx, key, token)
		if err != nil {
			t.Error("ack err", key, token, err)
		}
	}

	fmt.Println("batch dequeue", key, message)
}
