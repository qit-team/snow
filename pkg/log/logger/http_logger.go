package logger

import (
	"github.com/sirupsen/logrus"
	"context"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow/app/constants/common"
	"fmt"
)

var (
	hostname string
)

type withField struct {
	Key   string
	Value interface{}
}

//此结构的数据将会在挂靠到日志的一级键中体现
//demo: logger.Info(ctx, "curl", NewWithFiled("key1", "value1"), NewWithFiled("key2", "value2"), "msg1", "msg2")
func NewWithField(key string, value interface{}) *withField {
	return &withField{Key: key, Value: value}
}

//批量
func BatchNewWithField(data map[string]interface{}) (arr []*withField) {
	for k, v := range data {
		arr = append(arr, NewWithField(k, v))
	}
	return arr
}

func GetHostName() string {
	if hostname == "" {
		hostname, _ = os.Hostname()
		if hostname == "" {
			hostname = "unknown"
		}
	}
	return hostname
}

func formatLog(c context.Context, t string, args ...*withField) logrus.Fields {
	data := logrus.Fields{
		"type":     t,
		"host": GetHostName(),
	}

	if c != nil {
		traceId := GetTraceId(c)
		if traceId != "" {
			data["trace_id"] = traceId
		}

		switch c.(type) {
		case *gin.Context:
			ginC := c.(*gin.Context)
			data["domain"] = ginC.Request.Host
			data["sip"] = ginC.Request.RemoteAddr
			data["cip"] = ginC.ClientIP()
		}
	}

	for _, field := range args {
		if _, ok := data[field.Key]; !ok {
			data[field.Key] = field.Value
		}
	}

	return data
}

//获取traceId
func GetTraceId(c context.Context) string {
	traceId := c.Value("trace_id")
	if traceId != "" && traceId != nil {
		return fmt.Sprintf("%s", traceId)
	}
	switch c.(type) {
	case *gin.Context:
		return c.(*gin.Context).GetHeader(common.HeaderTraceId)
	}
	return ""
}

func Trace(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Trace(newMsg...)
}

func Debug(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Debug(newMsg...)
}

func Info(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Info(newMsg...)
}

func Warn(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Warn(newMsg...)
}

func Error(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Error(newMsg...)
}

func Fatal(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Fatal(newMsg...)
}

func Panic(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Panic(newMsg...)
}

//将日志消息分裂
func splitMsg(msg []interface{}) (withFields []*withField, newMsg []interface{}) {
	for _, v := range msg {
		switch v.(type) {
		case *withField:
			withFields = append(withFields, v.(*withField))
		default:
			newMsg = append(newMsg, v)
		}
	}
	return
}
