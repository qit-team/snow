package metric

import (
	"net/http"

	"github.com/qit-team/snow/app/utils/metric"

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
