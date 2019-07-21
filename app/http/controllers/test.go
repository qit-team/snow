package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow/app/http/formatters/bannerformatter"
	"github.com/qit-team/snow/app/services/bannerservice"
	"github.com/qit-team/snow/app/http/entities"
	"github.com/qit-team/snow/app/constants/errorcode"
	"time"
	"github.com/qit-team/snow-core/log/logger"
	"strconv"
)

//hello示例
func HandleHello(c *gin.Context) {
	logger.Debug(c, "hello", "test message")
	Success(c, "hello world!")
	return
}

//request和response的示例
func HandleTest(c *gin.Context) {
	request := new(entities.TestRequest)
	err := GenRequest(c, request)
	if err != nil {
		Error(c, errorcode.ParamError)
		return
	}

	response := new(entities.TestResponse)
	response.Name = request.Name
	response.Url = request.Url
	response.Id = time.Now().Unix()
	Success(c, response)
	return
}

//测试数据库服务示例
func GetBannerList(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	list, err := bannerservice.GetListByPid(1, limit, page)
	if err != nil {
		Error500(c)
		return
	}

	data := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"data":  bannerformatter.FormatList(list),
	}

	Success(c, data)
}
