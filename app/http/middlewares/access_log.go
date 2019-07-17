package middlewares

import (
	"github.com/qit-team/snow/pkg/log/accesslogger"
	"github.com/gin-gonic/gin"
	ginacl "github.com/hetiansu5/accesslog/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//忽略HEAD探针的日志
		if c.Request.Method != "HEAD" {
			ginacl.AccessLogFunc(accesslogger.GetAccessLogger())(c)
		}
	}
}
