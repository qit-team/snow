package middlewares

import (
	"github.com/qit-team/snow/pkg/log/accesslogger"
	"github.com/gin-gonic/gin"
	ginacl "github.com/hetiansu5/accesslog/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginacl.AccessLogFunc(accesslogger.GetAccessLogger())(c)
	}
}
