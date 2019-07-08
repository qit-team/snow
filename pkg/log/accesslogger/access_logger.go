package accesslogger

import (
	"github.com/qit-team/snow/pkg/log/logger"
	"github.com/hetiansu5/accesslog"
	coresio "github.com/hetiansu5/cores/io"
	"io"
)

func InitAccessLog(logHandler string, logDir string) (*accesslog.AccessLogger, error) {
	var writer io.Writer
	if logHandler == logger.HandlerStdout {
		writer = logger.GetStdOutWriter(logDir)
	} else {
		logFile := logDir + "/access.log"
		writerFile, err := coresio.NewRollingFileWriter(logFile, coresio.NewDailyRollingManager())
		if err != nil {
			return nil, err
		}
		writer = writerFile
	}

	acl, err := accesslog.NewLogger(accesslog.Output(writer), accesslog.Pattern(accesslog.JSONPattern))
	if err != nil {
		return nil, err
	}
	return acl, nil
}
