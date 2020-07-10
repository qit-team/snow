package routes

/**
 * 配置路由
 */
import (
	"github.com/qit-team/snow/app/http/controllers"
	"github.com/qit-team/snow/app/http/middlewares"
	"github.com/qit-team/snow/app/utils/metric"
	"github.com/qit-team/snow/config"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/http/middleware"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// api路由配置
func RegisterRoute(router *gin.Engine) {
	// middleware: 服务错误处理 => 生成请求id => access log
	router.Use(middlewares.ServerRecovery(), middleware.GenRequestId, middleware.GenContextKit, middleware.AccessLog())

	if config.GetConf().PrometheusCollectEnable && config.IsEnvEqual(config.ProdEnv) {
		router.Use(middlewares.CollectMetric())
		metric.Init(metric.EnableRuntime(), metric.EnableProcess())
		metricHandler := metric.Handler()
		router.GET("/metrics", func(ctx *gin.Context) {
			metricHandler.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	router.NoRoute(controllers.Error404)
	router.GET("/hello", controllers.HandleHello)
	router.POST("/test", controllers.HandleTest)
	router.POST("/test_validator", controllers.HandleTestValidator)

	// api版本
	v1 := router.Group("/v1")
	{
		v1.GET("/banner_list", controllers.GetBannerList)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
