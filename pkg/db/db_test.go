package db

import (
	"testing"
	"github.com/qit-team/snow/config"
	"github.com/go-xorm/xorm"
	"fmt"
	"time"
)

var engineGroup *xorm.EngineGroup

/**
 * Banner实体
 */
type Banner struct {
	Id        int64
	Pid       int
	Title     string
	ImageUrl  string    `xorm:"'img_url'"`
	Url       string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `xorm:"deleted"` //此特性会激发软删除
}

func init() {
	dbConf := config.DbConfig{
		Driver: "mysql",
		Master: config.DbBaseConfig{
			Host:     "127.0.0.1",
			User:     "root",
			Password: "123456",
			DBName:   "test",
		},
	}

	var err error
	engineGroup, err = NewEngineGroup(dbConf)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGet(t *testing.T) {
	banner := new(Banner)
	_, err := engineGroup.ID(1).Get(banner)
	if err != nil {
		t.Errorf("get error: %v", err)
	}
	fmt.Println(banner)
}
