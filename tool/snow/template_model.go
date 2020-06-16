package main

const (
	_tplModel = `package bannermodel

import (
	"sync"
	"time"

	"github.com/qit-team/snow-core/db"
)

var (
	once sync.Once
	m    *bannerModel
)
/**
 * Banner实体
 */
type Banner struct {
	Id        int64     ` + "`xorm:\"pk autoincr\"`" + ` //注：使用getOne 或者ID() 需要设置主键
	Pid       int
	Title     string
	ImageUrl  string    ` + "`xorm:\"'img_url'\"`" + `
	Url       string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time ` + "`xorm:\"deleted\"`" + ` //此特性会激发软删除
}

/**
 * 表名规则
 * @wiki http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-02/3.tags.html
 */
func (m *Banner) TableName() string {
	return "banner"
}

/**
 * 私有化，防止被外部new
 */
type bannerModel struct {
	db.Model //组合基础Model，集成基础Model的属性和方法
}

//单例模式
func GetInstance() *bannerModel {
	once.Do(func() {
		m = new(bannerModel)
		//m.DiName = "" //设置数据库实例连接，默认db.SingletonMain
	})
	return m
}

func (m *bannerModel) GetListByPid(pid int, limits ...int) (banners []*Banner, err error) {
	banners = make([]*Banner, 0)
	err = m.GetList(&banners, "pid = ?", []interface{}{pid}, limits)
	return
}
`

	_tplModelTest = `package bannermodel

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/db"
	"github.com/qit-team/snow-core/utils"
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
		t.Error("missing banner data")
	} else if banner.Id == 0 {
		t.Error("missing banner data")
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
`
)
