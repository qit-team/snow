package main

const (
	_tplControllerBase = `package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"{{.ModuleName}}/app/constants/errorcode"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/log/logger"
	"gopkg.in/go-playground/validator.v9"
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

type HTTPError struct {
	Code    int ` + "`json:\"code\" example:\"400\"`" + `
	Message string ` + "`json:\"message\" example:\"status bad request\"`" + `
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
	err = json.Unmarshal(body, request)
	if (err == nil) {
		validate := validator.New()
		errValidate := validate.Struct(request)
		if errValidate != nil {
			logger.Error(c, "param_validator_exception:" + c.Request.URL.Path, errValidate)
			return errValidate
		}
	}
	return err
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
	"strconv"
	"time"

	"{{.ModuleName}}/app/constants/errorcode"
	"{{.ModuleName}}/app/http/entities"
	"{{.ModuleName}}/app/http/formatters/bannerformatter"
	"{{.ModuleName}}/app/services/bannerservice"
	"{{.ModuleName}}/app/utils/httpclient"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/log/logger"
)

// hello示例
func HandleHello(c *gin.Context) {
	logger.Debug(c, "hello", "test message")
	client := httpclient.NewClient(c.Request.Context())
	resposne, err := client.R().Get("https://www.baidu.com")
	if err != nil {
		Error(c, errorcode.SystemError, err.Error())
		return
	}
	logger.Info(c, "HandleHello", resposne.String())
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

// 测试数据库服务示例
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
`

	_tplEntity = `package entities

//请求数据结构
type TestRequest struct {
	Name string ` + "`json:\"name\" example:\"snow\"`" + `
	Url  string ` + "`json:\"url\" example:\"github.com/qit-team/snow\"`" + `
}

//返回数据结构
type TestResponse struct {
	Id   int64  ` + "`json:\"id\" example:\"1\"`" + `
	Name string ` + "`json:\"name\" example:\"snow\"`" + `
	Url  string ` + "`json:\"url\" example:\"github.com/qit-team/snow\"`" + `
}

/*
 * validator.v9文档
 * 地址https://godoc.org/gopkg.in/go-playground/validator.v9
 * 列了几个大家可能会用到的，如有遗漏，请看上面文档
 */

//请求数据结构
type TestValidatorRequest struct {
	//tips，因为组件required不管是没传值或者传 0 or "" 都通过不了，但是如果用指针类型，那么0就是0，而nil无法通过校验
	Id   *int64 ` + "`json:\"id\" validate:\"required\" example:\"1\"`" + `
	Age  int ` + "`json:\"age\" validate:\"required,gte=0,lte=130\" example:\"20\"`" + `
	Name *string ` + "`json:\"name\" validate:\"required\" example:\"snow\"`" + `
	Email string ` + "`json:\"email\" validate:\"required,email\" example:\"snow@github.com\"`" + `
	Url  string ` + "`json:\"url\" validate:\"required\" example:\"github.com/qit-team/snow\"`" + `
	Mobile string ` + "`json:\"mobile\" validate:\"required\" example:\"snow\"`" + `
	RangeNum int ` + "`json:\"range_num\" validate:\"max=10,min=1\" example:\"3\"`" + `
	TestNum *int ` + "`json:\"test_num\" validate:\"required,oneof=5 7 9\" example:\"7\"`" + `
	Content *string ` + "`json:\"content\" example:\"snow\"`" + `
	Addresses []*Address ` + "`json:\"addresses\" validate:\"required,dive,required\"  `" + `
}

// Address houses a users address information
type Address struct {
	Street string ` + "`json:\"street\" validate:\"required\" example:\"huandaodonglu\"`" + `
	City   string ` + "`json:\"city\" validate:\"required\" example:\"xiamen\"`" + `
	Planet string ` + "`json:\"planet\" validate:\"required\" example:\"snow\"`" + `
	Phone  string ` + "`json:\"phone\" validate:\"required\" example:\"snow\"`" + `
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
	"testing"

	"{{.ModuleName}}/app/models/bannermodel"
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

	_tplMetric = `package metric

import (
	"net/http"

	"{{.ModuleName}}/app/utils/metric"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	HOST   = "host"
	PATH   = "path"   // 路径
	METHOD = "method" // 方法
	CODE   = "code"   // 错误码

	// metric
	ALL_REQ_TOTAL_COUNT = "all_req_total_count" // 所有URL总请求数
	ALL_REQ_COST_TIME   = "all_req_cost_time"   // 所有URL请求耗时

	REQ_TOTAL_COUNT = "req_total_count" // 每个URL总请求数
	REQ_COST_TIME   = "req_cost_time"   // 每个URL请求耗时
)

func init() {
	metric.RegisterCollector(reqTotalCounter, reqCostTimeObserver, allReqTotalCounter, allReqCostTimeObserver)
}

var (
	reqTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: REQ_TOTAL_COUNT,
	}, []string{PATH, METHOD})

	reqCostTimeObserver = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: REQ_COST_TIME,
		Buckets: []float64{
			100,
			200,
			500,
			1000,
			3000,
			5000,
		},
	}, []string{PATH, METHOD})

	allReqTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: ALL_REQ_TOTAL_COUNT,
	}, []string{HOST})

	allReqCostTimeObserver = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: ALL_REQ_COST_TIME,
		Buckets: []float64{
			100,
			200,
			500,
			1000,
			3000,
			5000,
		},
	}, []string{HOST})
)

func AddReqCount(req *http.Request) {
	reqTotalCounter.WithLabelValues(req.URL.Path, req.Method).Inc()
}

func CollectReqCostTime(req *http.Request, ms int64) {
	reqCostTimeObserver.WithLabelValues(req.URL.Path, req.Method).Observe(float64(ms))
}

func AddAllReqCount(req *http.Request) {
	allReqTotalCounter.WithLabelValues(req.Host).Inc()
}

func CollectAllReqCostTime(req *http.Request, ms int64) {
	allReqCostTimeObserver.WithLabelValues(req.Host).Observe(float64(ms))
}
`

	_tplSkyWalkingTracer = `package trace

import (
	"sync"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"{{.ModuleName}}/config"
)

var (
	tracer *go2sky.Tracer
	lock   sync.Mutex
)

func Tracer() (*go2sky.Tracer, error) {
	if tracer == nil {
		lock.Lock()
		defer lock.Unlock()
		if tracer == nil {
			err := InitTracer(config.GetConf().ServiceName, config.GetConf().SkyWalkingOapServer)
			if err != nil {
				return nil, err
			}
		}
	}
	return tracer, nil
}

func InitTracer(serviceName, skyWalkingOapServer string) error {
	var (
		report go2sky.Reporter
		err    error
	)
	report, err = reporter.NewGRPCReporter(skyWalkingOapServer)
	if err != nil {
		return err
	}
	tracer, err = go2sky.NewTracer(serviceName, go2sky.WithReporter(report))
	if err != nil {
		return err
	}
	return nil
}
`

	_tplMiddleWare = `package middlewares

import (
	"encoding/json"
	syslog "log"
	"net/http/httputil"
	"runtime/debug"

	"{{.ModuleName}}/app/constants/logtype"
	"{{.ModuleName}}/config"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/log/logger"
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

	_tplMiddleWreMetric = `package middlewares

import (
	"time"

	"{{.ModuleName}}/app/http/metric"

	"github.com/gin-gonic/gin"
)

func CollectMetric() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		dur := time.Now().Sub(start).Milliseconds()

		metric.AddAllReqCount(ctx.Request)
		metric.CollectAllReqCostTime(ctx.Request, dur)
		metric.AddReqCount(ctx.Request)
		metric.CollectReqCostTime(ctx.Request, dur)
	}
}
`

	_tplMiddleWreSkyWalkingTracer = `package middlewares

import (
	"fmt"
	"strconv"
	"time"

	"{{.ModuleName}}/app/http/trace"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/log/logger"
)

const (
	componentIDGOHttpServer = 5004
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer, err := trace.Tracer()
		if err != nil {
			logger.Error(c, "Trace", err.Error())
			c.Next()
			return
		}
		r := c.Request
		operationName := fmt.Sprintf("/%s%s", r.Method, r.URL.Path)
		span, ctx, err := tracer.CreateEntrySpan(c, operationName, func() (string, error) {
			return r.Header.Get(propagation.Header), nil
		})
		if err != nil {
			logger.Error(c, "Trace", err.Error())
			c.Next()
			return
		}
		span.SetComponent(componentIDGOHttpServer)
		span.Tag(go2sky.TagHTTPMethod, r.Method)
		span.Tag(go2sky.TagURL, fmt.Sprintf("%s%s", r.Host, r.URL.Path))
		span.SetSpanLayer(v3.SpanLayer_Http)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		code := c.Writer.Status()
		if code >= 400 {
			span.Error(time.Now(), fmt.Sprintf("Error on handling request, statusCode: %d", code))
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(code))
		span.End()
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
	"{{.ModuleName}}/app/http/trace"
	"{{.ModuleName}}/app/utils/metric"
	"{{.ModuleName}}/config"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/http/middleware"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//api路由配置
func RegisterRoute(router *gin.Engine) {
	//middleware: 服务错误处理 => 生成请求id => access log
	router.Use(middlewares.ServerRecovery(), middleware.GenRequestId, middleware.GenContextKit, middleware.AccessLog())

	if config.GetConf().PrometheusCollectEnable && config.IsEnvEqual(config.ProdEnv) {
		router.Use(middlewares.CollectMetric())
		metric.Init(metric.EnableRuntime(), metric.EnableProcess())
		metricHandler := metric.Handler()
		router.GET("/metrics", func(ctx *gin.Context) {
			metricHandler.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	if len(config.GetConf().SkyWalkingOapServer) > 0 && config.IsEnvEqual(config.ProdEnv) {
		err := trace.InitTracer(config.GetConf().ServiceName, config.GetConf().SkyWalkingOapServer)
		if err != nil {
			logger.Error(nil, "InitTracer", err.Error())
		} else {
			router.Use(middlewares.Trace())
		}
	}

	router.NoRoute(controllers.Error404)
	router.GET("/hello", controllers.HandleHello)
	router.POST("/test", controllers.HandleTest)
    router.POST("/test_validator", controllers.HandleTestValidator)
	
    //api版本
	v1 := router.Group("/v1")
	{
		v1.GET("/banner_list", controllers.GetBannerList)
	}
    
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
`
)
