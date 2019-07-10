package db

import (
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq" //postgres
	//_ "github.com/mattn/go-sqlite3" //sqlite3
	//_ "github.com/denisenkom/go-mssqldb" //mssql
	"github.com/qit-team/snow/config"
	"fmt"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"errors"
)

const (
	defaultTimeout = 10
	defaultCharset = "utf8mb4"
)

func NewEngineGroup(dbConf config.DbConfig) (*xorm.EngineGroup, error) {
	master, err := newConn(dbConf.Driver, dbConf.Master, dbConf.Option)
	if err != nil {
		panicConnectionErr(dbConf.Driver, dbConf.Master.Host, dbConf.Master.Port, err)
	}

	slaves := make([]*xorm.Engine, len(dbConf.Slaves))
	for k, slaveConf := range dbConf.Slaves {
		slave, err := newConn(dbConf.Driver, slaveConf, dbConf.Option)
		if err != nil {
			panicConnectionErr(dbConf.Driver, slaveConf.Host, slaveConf.Port, err)
		}
		slaves[k] = slave
	}

	return xorm.NewEngineGroup(master, slaves)
}

func newConn(driver string, base config.DbBaseConfig, option config.DbOptionConfig) (db *xorm.Engine, err error) {
	dsn := formatDSN(driver, base, option)
	if dsn == "" {
		return nil, errors.New(fmt.Sprintf("missing db driver %s or db config", driver))
	}
	db, err = xorm.NewEngine(driver, dsn)
	if err != nil {
		return
	}
	
	//设置表名和字段的映射规则：驼峰转下划线
	db.SetMapper(core.SnakeMapper{})

	//设置资源池等配置
	if option.MaxIdle > 0 {
		db.SetMaxIdleConns(option.MaxIdle)
	}
	if option.MaxConns > 0 {
		db.SetMaxOpenConns(option.MaxConns)
	}
	if option.IdleTimeout > 0 {
		db.SetConnMaxLifetime(time.Second * option.IdleTimeout)
	}
	return
}

/**
 * 各驱动的dsn
 * @wiki http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/
 */
func formatDSN(driver string, base config.DbBaseConfig, option config.DbOptionConfig) string {
	switch driver {
	case "mysql":
		return formatMysqlDSN(base, option)
	case "postgres":
		return formatPostgresDSN(base, option)
	case "sqlite3":
		return formatSqlite3DSN(base, option)
	case "mssql":
		return formatMssqlDSN(base, option)
	}
	return ""
}

//Mysql DSN
func formatMysqlDSN(base config.DbBaseConfig, option config.DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 3306)
	charset := option.Charset
	if charset == "" {
		charset = defaultCharset
	}
	timeout := option.ConnectTimeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=%s&parseTime=true&loc=Local",
		base.User, base.Password, base.Host, port, base.DBName, option.ConnectTimeout, option.Charset)
}

//PostgreSQL DSN
func formatPostgresDSN(base config.DbBaseConfig, option config.DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 5432)
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		base.Host, port, base.User, base.DBName, base.Password)
}

//qlite3 DSN
func formatSqlite3DSN(base config.DbBaseConfig, option config.DbOptionConfig) string {
	return base.DBName
}

//SQL Server DSN
func formatMssqlDSN(base config.DbBaseConfig, option config.DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 1433)
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		base.User, base.Password, base.Host, port, base.DBName)
}

func getPortOrDefault(port int, defaultPort int) int {
	if port == 0 {
		return defaultPort
	}
	return port
}

func panicConnectionErr(driver string, host string, port int, err error) {
	panic(fmt.Sprintf("%s connect error %s:%d, error:%v", driver, host, port, err))
}
