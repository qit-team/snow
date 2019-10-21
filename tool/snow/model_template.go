package main

const (
	tplModel = `package {{.PackageName}}

import (
	"github.com/qit-team/snow-core/db"
	"sync"
	"time"
)

var (
	once sync.Once
	m    *{{.ModelName}}
)

//实体
{{.TableEntity}}

//表名
func (m *{{.Entity}}) TableName() string {
	return "{{.Table}}"
}

//私有化，防止被外部new
type {{.ModelName}} struct {
	db.Model //组合基础Model，集成基础Model的属性和方法
}

//单例模式
func GetInstance() *{{.ModelName}} {
	once.Do(func() {
		m = new({{.ModelName}})
		//m.DiName = "" //设置数据库实例连接，默认db.SingletonMain
	})
	return m
}
`
)
