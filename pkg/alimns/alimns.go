package alimns

import (
	"github.com/qit-team/snow/config"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"fmt"
	"errors"
)

//依赖注入用的函数
func NewMnsClient(mnsConfig config.MnsConfig) (client ali_mns.MNSClient, err error) {
	//2.1初始化mns client
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			fmt.Println("ali_mns client init error", err)
		}
	}()

	if mnsConfig.Url != "" {
		client = ali_mns.NewAliMNSClient(mnsConfig.Url,
			mnsConfig.AccessKeyId,
			mnsConfig.AccessKeySecret)
	}
	return
}

func GetMnsBasicQueue(client ali_mns.MNSClient, queueName string) ali_mns.AliMNSQueue {
	var defaultQueue ali_mns.AliMNSQueue

	//根据client创建manager
	queueManager := ali_mns.NewMNSQueueManager(client)
	err := queueManager.CreateQueue(queueName, 0, 65536, 345600, 30, 0, 3)
	if err != nil && !ali_mns.ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
		fmt.Println(err)
		return defaultQueue
	}
	//最终的最小执行单元queue
	return ali_mns.NewMNSQueue(queueName, client)
}
