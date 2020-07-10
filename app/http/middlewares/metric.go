package middlewares

import (
	"time"

	"github.com/qit-team/snow/app/http/metric"

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
