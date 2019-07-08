package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
)

//app.log_handler为file时，日志格式为:[time(ISO8601)]  [host]  [type(service.module.function)]  [req_id]  [server_ip]  [client_ip]  [message(json:code,message,file,line,trace,biz_data)]
//app.log_handler为stdout时，日志格式为:{"t": "time(ISO8601)", "lvl": "level", "h": "host", "type": "type(service.module.function)", "reqid": "req_id", "sip": "server_ip", "cip": "client_ip", "msg": {"code": 0, "message": "xxx", "file": "file", "line": 0}}

const HandlerFile = "file"
const HandlerStdout = "stdout"

func GetStdOutWriter(path string) (writer *os.File) {
	//此处命名管道会阻塞，直到有进程读取了这个命名管道
	writer, err := os.OpenFile(path, os.O_WRONLY, 777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file, %v\n", err)
	}
	return
}

func InitLog(logHandler string, logDir string, logLevel string) (*logrus.Logger, error) {
	logger := logrus.New()

	//设置日志等级
	level, err := logrus.ParseLevel(logLevel)
	if err == nil {
		logger.SetLevel(level)
	}

	//设置日志输出格式
	logger.Formatter = &logrus.JSONFormatter{}

	//设置日志输出方式 标准输出或文件
	if logHandler == HandlerStdout {
		writer := GetStdOutWriter(logDir)
		logger.SetOutput(writer)
	} else {
		rollHook, err := NewRollHook(logger, logDir, "snow")
		if err != nil {
			return nil, err
		}
		logger.Hooks.Add(rollHook)
	}

	return logger, nil
}
