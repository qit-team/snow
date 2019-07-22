package main

const (
	_tplControllerBase = `package controllers

import (
	"{{.ModuleName}}/app/constants/errorcode"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

/**
 * 成功时返回
 */
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":        errorcode.Success,
		"message":     "ok",
		"request_uri": c.Request.URL.Path,
		"data":        data,
	})
	c.Abort()
}

/**
 * 失败时返回
 */
func Error(c *gin.Context, code int, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = errorcode.GetMsg(code)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":        code,
		"message":     message,
		"request_uri": c.Request.URL.Path,
		"data":        make(map[string]string),
	})
	c.Abort()
}

func Error404(c *gin.Context) {
	Error(c, errorcode.NotFound, "路由不存在")
}

func Error500(c *gin.Context) {
	Error(c, errorcode.SystemError)
}

/**
 * 将请求的body转换为request数据结构
 * @param c
 * @param request  传入request数据结构的指针 如 new(TestRequest)
 */
func GenRequest(c *gin.Context, request interface{}) (err error) {
	body, err := ReadBody(c)
	if err != nil {
		return
	}
	return json.Unmarshal(body, request)
}

//重复读取body
func ReadBody(c *gin.Context) (body []byte, err error) {
	body, err = ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return
}
`

	_tplControllerTest = `package controllers

import (
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/app/http/formatters/bannerformatter"
	"{{.ModuleName}}/app/services/bannerservice"
	"{{.ModuleName}}/app/http/entities"
	"{{.ModuleName}}/app/constants/errorcode"
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
`

	_tplEntity = `package entities

//请求数据结构
type TestRequest struct {
	Name string ` + "`json:\"name\"`" + `
	Url  string ` + "`json:\"url\"`" + `
}

//返回数据结构
type TestResponse struct {
	Id   int64  ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
	Url  string ` + "`json:\"url\"`" + `
}`

	_tplFormatter = `package bannerformatter

import (
	"{{.ModuleName}}/app/models/bannermodel"
)

type BannerFormatter struct {
	Id    int    ` + "`json:\"id\"`" + `
	Title string ` + "`json:\"title\"`" + `
	Img   string ` + "`json:\"image\"`" + `
	Url   string ` + "`json:\"url\"`" + `
}

func FormatList(bannerList []*bannermodel.Banner) (res []*BannerFormatter) {
	res = make([]*BannerFormatter, len(bannerList))

	for k, banner := range bannerList {
		one := FormatOne(banner)
		res[k] = one
	}

	return res
}

//单条消息的格式化，
func FormatOne(banner *bannermodel.Banner) (res *BannerFormatter) {
	res = &BannerFormatter{
		Id:    int(banner.Id),
		Title: banner.Title,
		Img:   banner.ImageUrl,
		Url:   banner.Url,
	}
	return
}`

	_tplFormatterTest = `package bannerformatter

import (
	"{{.ModuleName}}/app/models/bannermodel"
	"testing"
)

func TesFormatOne(t *testing.T) {
	a := &bannermodel.Banner{
		Id:       1,
		Title:    "test",
		ImageUrl: "http://x/1.jpg",
		Url:      "http://x",
		Status:   "1",
	}
	b := FormatOne(a)
	if b.Title != a.Title || b.Img != a.ImageUrl || b.Url != a.Url {
		t.Error("FormatOne not same")
	}
}

func TesFormatList(t *testing.T) {
	a := make([]*bannermodel.Banner, 2)
	a[0] = &bannermodel.Banner{
		Id:       1,
		Title:    "test",
		ImageUrl: "http://x1/1.jpg",
		Url:      "http://x1",
		Status:   "1",
	}
	a[1] = &bannermodel.Banner{
		Id:       2,
		Title:    "test2",
		ImageUrl: "http://x/2.jpg",
		Url:      "http://x2",
		Status:   "2",
	}
	b := FormatList(a)
	for k, v := range b {
		if v.Title != a[k].Title || v.Img != a[k].ImageUrl || v.Url != a[k].Url {
			t.Error("FormatList not same")
		}
	}
}
`

	_tplMiddleWare = `package middlewares

import (
	"encoding/json"
	"{{.ModuleName}}/app/constants/logtype"
	"{{.ModuleName}}/config"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/gin-gonic/gin"
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
					//本地开发 debug 模式开启时输出错误信息到shell
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

		//before request

		c.Next()

		//after request
	}
}
`

	_tplRoute = `package routes

/**
 * 配置路由
 */
import (
	"{{.ModuleName}}/app/http/controllers"
	"{{.ModuleName}}/app/http/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/http/middleware"
)

//api路由配置
func RegisterRoute(router *gin.Engine) {
	//middleware: 服务错误处理 => 生成请求id => access log
	router.Use(middlewares.ServerRecovery(), middleware.GenRequestId, middleware.GenContextKit, middleware.AccessLog())

	router.NoRoute(controllers.Error404)
	router.GET("/hello", controllers.HandleHello)
	router.POST("/test", controllers.HandleTest)

	//api版本
	v1 := router.Group("/v1")
	{
		v1.GET("/banner_list", controllers.GetBannerList)
	}
}
`
)
