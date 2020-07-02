package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

const (
	defaultDSN        = "root:123456@tcp(localhost:3306)/test"
	defaultDriverType = "mysql"
)

var (
	m  *model
	mt *modelTemplate
)

type field struct {
	Name string
	Type string
	Tag  string
}

type model struct {
	Name       string //模型名称
	Path       string //项目目录
	Table      string //表名，默认同模型名称
	DSN        string //数据库连接的dsn配置
	DB         string //数据库名称
	DriverType string //驱动类型，默认mysql
}

type modelTemplate struct {
	SnakeName   string //模型名称的蛇型法
	PackageName string //模型名称的包名
	ModelName   string //模型名称的驼峰法
	Entity      string //表的实体结构名称
	Table       string //表名
	TableEntity string //表的结构定义
}

func (m *model) getDriverType() string {
	if m.DriverType == "" {
		return defaultDriverType
	}
	return m.DriverType
}

func (m *model) getDSN() string {
	var dsn string
	if m.DSN != "" {
		dsn = m.DSN
	} else {
		dsn = os.Getenv("SNOW_DSN")
	}
	if dsn == "" {
		dsn = defaultDSN
	}
	if m.DB != "" {
		arr := strings.SplitN(dsn, "/", 2)
		if len(arr) < 2 {
			arr = append(arr, m.DB)
		} else {
			arr[1] = m.DB
		}
		dsn = strings.Join(arr, "/")
	}
	return dsn
}

func init() {
	m = new(model)
	mt = new(modelTemplate)
}

//new model
func runModel(ctx *cli.Context) (err error) {
	if ctx.Args().Len() == 0 {
		return errors.New("required model name")
	}

	m.Name = ctx.Args().First()
	if m.Name == "" {
		return errors.New("model name is empty")
	}

	if m.Table == "" {
		m.Table = m.Name
	}

	if m.Path == "" {
		m.Path, _ = os.Getwd()
	} else {
		if !isDirExist(m.Path) {
			return errors.New("project directory is not exist")
		}
	}

	snakeMapper := new(core.SnakeMapper)
	mt.SnakeName = snakeMapper.Obj2Table(m.Name)
	mt.PackageName = packageCase(mt.SnakeName) + "model"
	mt.Entity = snakeMapper.Table2Obj(mt.SnakeName)
	mt.ModelName = snakeMapper.Table2Obj(mt.SnakeName) + "Model"
	mt.Table = m.Table

	//create model directory
	path := strings.Join([]string{m.Path, "app/models", mt.PackageName}, string(os.PathSeparator))
	if err = os.MkdirAll(path, 0755); err != nil {
		return
	}

	//dsn
	dsn := m.getDSN()
	if dsn == "" {
		return errors.New("dsn is empty")
	}

	//连接数据库
	engine, err := xorm.NewEngine(m.getDriverType(), dsn)
	if err != nil {
		return
	}

	//获取数据库表定义
	tables, err := engine.DBMetas()
	if err != nil {
		return
	}

	//获取表结构定义
	var ok bool
	for _, table := range tables {
		if table.Name != mt.Table {
			continue
		}
		ok = true
		mt.TableEntity = genTableEntity(table)
		break
	}

	if !ok {
		return errors.New("cannot find related table")
	}

	//将模板写入文件
	file := path + string(os.PathSeparator) + mt.SnakeName + ".go"
	err = write(file, tplModel, mt)
	if err != nil {
		return
	}

	//输出提示信息
	fmt.Printf("Table: %s\n", mt.Table)
	fmt.Printf("Entity: %s\n", mt.Entity)
	fmt.Printf("Model Name: %s\n", mt.ModelName)
	fmt.Printf("Package Name: %s\n", mt.PackageName)
	fmt.Printf("Directory: %s\n\n", path)
	fmt.Println("The model has been created.")
	return nil
}

//将蛇形命名法转换为go包名连写命名法
func packageCase(name string) string {
	name = strings.ToLower(name)
	rs := []rune(name)
	var buffer bytes.Buffer
	var s string
	for _, v := range rs {
		s = string(v)
		if s == "_" {
			continue
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

//目录是否存在
func isDirExist(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

//将数据库table的定义转换为数据结构实体定义
func genTableEntity(table *schemas.Table) string {
	columns := table.Columns()
	snakeMapper := new(core.SnakeMapper)

	RowInfoList := make([]*field, 0)
	for _, column := range columns {
		sqlType := column.SQLType

		rowInfo := &field{
			// 变量名
			Name: snakeMapper.Table2Obj(column.Name),
		}

		// 变量类型
		if sqlType.Name == core.DateTime {
			rowInfo.Type = "time.Time"
		} else if sqlType.Name == core.TimeStamp {
			rowInfo.Type = "time.Time"
		} else {
			rowInfo.Type = schemas.SQLType2Type(sqlType).Name()
		}

		// xorm注释
		tag := "`xorm:\"'" + column.Name + "'"
		if sqlType.Name == core.DateTime {
			tag += " datetime"
		} else if sqlType.Name == core.TimeStamp {
			tag += " timestamp"
		} else if sqlType.Name == core.BigInt {
			tag += " bigint(" + strconv.Itoa(sqlType.DefaultLength) + ")"
		} else if sqlType.Name == core.Int {
			tag += " int(" + strconv.Itoa(sqlType.DefaultLength) + ")"
		} else if sqlType.Name == core.Decimal {
			tag += " decimal(" + strconv.Itoa(sqlType.DefaultLength) + "," + strconv.Itoa(sqlType.DefaultLength2) + ")"
		} else if sqlType.Name == core.Varchar {
			tag += " varchar(" + strconv.Itoa(sqlType.DefaultLength) + ")"
		} else if sqlType.Name == core.Char {
			tag += " char(" + strconv.Itoa(sqlType.DefaultLength) + ")"
		} else {
			tag += " " + sqlType.Name
		}
		//特殊字段加点盐
		if column.Name == "id" {
			tag += " pk autoincr"
		} else if column.Name == "deleted_at" {
			tag += " deleted"
		}
		tag += "\"`"

		rowInfo.Tag = tag
		RowInfoList = append(RowInfoList, rowInfo)
	}

	maxVarNameStrLen := 0
	maxVarTypeStrLen := 0
	for _, field := range RowInfoList {
		if len(field.Name) > maxVarNameStrLen {
			maxVarNameStrLen = len(field.Name)
		}

		if len(field.Type) > maxVarTypeStrLen {
			maxVarTypeStrLen = len(field.Type)
		}
	}
	maxVarNameStrLen += 1
	maxVarTypeStrLen += 1

	// 生成文件内容
	structText := "type " + mt.Entity + " struct {\n"
	for _, field := range RowInfoList {
		structText += "\t"
		structText += field.Name
		if len(field.Name) < maxVarNameStrLen {
			padSpanTotal := maxVarNameStrLen - len(field.Name)
			for index := 0; index < padSpanTotal; index++ {
				structText += " "
			}
		}

		structText += field.Type
		if len(field.Type) < maxVarTypeStrLen {
			padSpanTotal := maxVarTypeStrLen - len(field.Type)
			for index := 0; index < padSpanTotal; index++ {
				structText += " "
			}
		}
		structText += field.Tag
		structText += "\n"
	}
	structText += `}`
	return structText
}
