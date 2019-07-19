package models

import (
	"github.com/qit-team/snow/pkg/db"
	"github.com/go-xorm/xorm"
	"errors"
)

var (
	ErrIdsEmpty = errors.New("ids is empty")
)

/**
 * 基础model
 */
type Model struct {
	DiName string //依赖注入的别名
}

/**
 * 获取数据库实例
 * @wiki http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-02/4.columns.html
 */
func (m *Model) GetDb(args ...string) *xorm.EngineGroup {
	if len(args) > 0 {
		return db.GetDb(args[0])
	} else if m.DiName != "" {
		return db.GetDb(m.DiName)
	} else {
		return db.GetDb()
	}
}

/**
 * 查询主键ID的记录
 * @param id 主键ID
 * @param bean 数据结构实体
 * @return has 是否有记录
 */
func (m *Model) GetOne(id interface{}, bean interface{}) (has bool, err error) {
	return m.GetDb().ID(id).Get(bean)
}

/**
 * 查询多个主键ID的记录
 * @param ids 主键ID分片
 * @param beans 数据结构实体分片
 */
func (m *Model) GetMulti(ids []interface{}, beans interface{}) error {
	if len(ids) == 0 {
		return ErrIdsEmpty
	}
	return m.GetDb().In("id", ids...).Find(beans)
}

/**
 * 插入记录
 * @param beans... 可支持插入连续多个记录
 */
func (m *Model) Insert(beans ...interface{}) (int64, error) {
	return m.GetDb().Insert(beans...)
}

/**
 * 更新某个主键ID的数据
 * @param id 主键ID
 * @param bean 数据结构实体
 * @param mustColumns... 因为默认Update只更新非0，非”“，非bool的字段，需要配合此字段
 * @param
 */
func (m *Model) Update(id interface{}, bean interface{}, mustColumns ...string) (int64, error) {
	if len(mustColumns) > 0 {
		return m.GetDb().MustCols(mustColumns...).ID(id).Update(bean)
	} else {
		return m.GetDb().ID(id).Update(bean)
	}
}

/**
 * 删除单个记录 -- 如果有开启delete特性，会触发软删除
 * @param id 主键ID
 * @param bean 数据结构实体
 */
func (m *Model) Delete(id interface{}, bean interface{}) (int64, error) {
	return m.GetDb().ID(id).Delete(bean)
}

/**
 * 查询多个主键ID的记录
 * @param ids 主键ID分片
 * @param bean 数据结构实体
 */
func (m *Model) DeleteMulti(ids []interface{}, bean interface{}) (int64, error) {
	if len(ids) == 0 {
		return 0, ErrIdsEmpty
	}
	return m.GetDb().In("id", ids...).Delete(bean)
}

/**
 * 查询多个主键ID的记录
 * @param beans 数据结构实体分片 eg. &banners 其中 banners := make([]*Banner, 0)
 * @params sql  eg. "age > ? or name = ?"
 * @params values eg. []interfaces{}{30, "hts"}
 * @Param []int limit 可选 eg. []int{} 不限量 []int{30} 前30个 []int{30, 20} 从第20个后的前30个
 * @param string order 可选 eg.  "id desc" 单个 "uid desc,status asc" 多个
 */
func (m *Model) GetList(beans interface{}, sql string, values []interface{}, args ...interface{}) (err error) {
	if len(args) > 0 {
		var (
			order string
			limit int
			start int
		)

		limits, ok := args[0].([]int)
		if ok && len(limits) > 0 {
			limit = limits[0]
			if len(limits) > 1 {
				start = limits[1]
			}
		}

		if len(args) > 1 {
			order, _ = args[1].(string)
		}

		return m.GetDb().Where(sql, values...).OrderBy(order).Limit(limit, start).Find(beans)
	} else {
		return m.GetDb().Where(sql, values...).Find(beans)
	}
}
