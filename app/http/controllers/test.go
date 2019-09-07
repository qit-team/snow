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

// hello示例
func HandleHello(c *gin.Context) {
	logger.Debug(c, "hello", "test message")
	Success(c, "hello world!")
	return
}

// request和response的示例
// HandleTest godoc
// @Summary request和response的示例
// @Description request和response的示例
// @Tags snow
// @Accept  json
// @Produce  json
// @Param test body entities.TestRequest true "test request"
// @Success 200 {array} entities.TestResponse
// @Failure 400 {object} controllers.HTTPError
// @Failure 404 {object} controllers.HTTPError
// @Failure 500 {object} controllers.HTTPError
// @Router /test [post]
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


// validator的示例
// HandleTestValidator godoc
// @Summary HandleTestValidator的示例
// @Description HandleTestValidator的示例
// @Tags snow
// @Accept  json
// @Produce json
// @Param testValidator body entities.TestValidatorRequest true "example of validator"
// @Success 200 {array} entities.TestValidatorRequest
// @Failure 400 {object} controllers.HTTPError
// @Failure 404 {object} controllers.HTTPError
// @Failure 500 {object} controllers.HTTPError
// @Router /test_validator [post]
func HandleTestValidator(c *gin.Context) {
	request := new(entities.TestValidatorRequest)
	err := GenRequest(c, request)
	if err != nil {
		Error(c, errorcode.ParamError)
		return
	}

	Success(c, request)
	return
}
