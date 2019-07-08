package redis

import (
	"fmt"
	"testing"
	"reflect"
	"strconv"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow/config"
)

var client *redis_pool.ReplicaPool

func init() {
	conf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	var err error
	client, err = NewRedisClient(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_Set(t *testing.T) {
	value := 11
	res, _ := client.Set("hts", value)
	t.Log(res, reflect.TypeOf(res))
	if res == false {
		t.Error("set error")
	}

	res1, _ := client.Get("hts")
	t.Log(res1, reflect.TypeOf(res1))
	if res1 == "" {
		t.Error("get error")
	} else if res1 != strconv.Itoa(value) {
		t.Error("not same")
	}
}
