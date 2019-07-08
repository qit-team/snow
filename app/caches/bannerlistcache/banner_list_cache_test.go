package bannerlistcache

import (
	"context"
	"fmt"
	"github.com/qit-team/snow/config"
	"github.com/qit-team/snow/pkg/redis"
	"testing"
)

func init() {
	//加载配置文件
	conf, err := config.Load("../../../.env")
	if err != nil {
		fmt.Println(err)
	}

	//注册redis类
	err = (&redis.Provider{}).Register(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_GetMulti(t *testing.T) {
	ctx := context.TODO()
	cache := GetCache()
	res, _ := cache.Set(ctx, "1000", "a", cache.TTL)
	if res != true {
		t.Errorf("set key:%s is error", "1000")
	}

	keys := []string{"1000", "-2000", "9999"}
	cacheList, err := cache.GetMulti(ctx, keys...)
	if err != nil {
		t.Errorf("getMulti error:%s", err.Error())
	}
	fmt.Println(cacheList)
	i := 0
	for k, v := range cacheList {
		i++
		if k == "1000" {
			if v != "a" {
				t.Errorf("value of key:%s is error %v", k, v)
			}
		} else {
			if v != "" {
				t.Errorf("value of key:%s is error %v", k, v)
			}
		}
	}
	if i != len(keys) {
		t.Errorf("count of cache key is error: %d", i)
	}
}
