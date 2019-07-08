package middlewares

import (
	"github.com/qit-team/snow/app/constants/errorcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, code int, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = errorcode.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"msg":         message,
		"request_uri": c.Request.URL.Path,
		"data":        make(map[string]string),
	})
	c.Abort()
}

//所有form表单的请求数据
func getFormData(c *gin.Context) map[string]interface{} {
	if c.Request.Form == nil {
		c.Request.ParseMultipartForm(32 << 20)
	}

	data := make(map[string]interface{})
	for k, v := range c.Request.Form {
		data[k] = v[0]
	}
	return data
}
