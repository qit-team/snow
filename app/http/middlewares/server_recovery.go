package middlewares

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow/app/constants/logtype"
	"github.com/qit-team/snow/config"
	syslog "log"
	"net/http/httputil"
	"runtime/debug"
)

func ServerRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				msg := map[string]interface{}{
					"error":   err,
					"request": string(httpRequest),
					"stack":   string(debug.Stack()),
				}
				msgJson, _ := json.Marshal(msg)
				logger.GetLogger().Error(string(msgJson), logtype.GoPanic, c)

				if config.IsDebug() {
					// 本地开发 debug 模式开启时输出错误信息到shell
					syslog.Println(err)
				}

				c.JSON(500, gin.H{
					"code":        500,
					"msg":         "system error",
					"request_uri": c.Request.URL.Path,
					"data":        make(map[string]string),
				})
			}
		}()

		// before request

		c.Next()

		// after request
	}
}
