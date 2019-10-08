package bannermodel

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/db"
	"github.com/qit-team/snow-core/utils"
	"testing"
)

func init() {
	m := config.DbBaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "test",
	}
	dbConf := config.DbConfig{
		Driver: "mysql",
		Master: m,
	}

	err := db.Pr.Register("db", dbConf, true)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetOne(t *testing.T) {
	bannerModel := GetInstance()
	banner := new(Banner)
	res, err := bannerModel.GetOne(1, banner)
	if err != nil {
		t.Error(err)
	} else if res != true {
		t.Error("missing bannner data")
	} else if banner.Id == 0 {
		t.Error("missing bannner data")
	}
	fmt.Println(utils.JsonEncode(banner))
}

func TestGetList(t *testing.T) {
	bannerModel := GetInstance()
	banners := make([]*Banner, 0)
	err := bannerModel.GetList(&banners, "pid >= ?", []interface{}{1}, []int{10}, "status desc, id desc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(utils.JsonEncode(banners))
}
