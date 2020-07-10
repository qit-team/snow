package metric

// prometheus metric：unique identifier: name and optional key-value pairs called labels
//   1. name regexp: [a-zA-Z_:][a-zA-Z0-9_:]*
//   2. label name regexp: [a-zA-Z_][a-zA-Z0-9_]*
//   3. Label names beginning with __ are reserved for internal use.
//   4. Label values may contain any Unicode characters.
//   5. notation: <metric name>{<label name>=<label value>, ...}
//      for example: api_http_requests_total{method="POST", handler="/messages"}
// A label with an empty label value is considered equivalent to a label that does not exist.

// each sample consists of :
//   - a float64 value
//   - a millisecond-precision timestamp

// metric type:
//   - Counter
//      A cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart.
//   - Gauge
//      A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
//   - Histogram
//      A histogram samples observations (usually things like request durations or response sizes) and counts them in configurable buckets. It also provides a sum of all observed values.
//   - Summary
//      Similar to a histogram, a summary samples observations (usually things like request durations and response sizes). While it also provides a total count of observations and a sum of all observed values, it calculates configurable quantiles over a sliding time window.
// metric:
//   Counter:
//      - req_total_count
//      - req_failed_count
//   Gauge:
//      - heap_inuse_size
//      - heap_total_size
//      - heap_object_num
//      - goroutine_num
//   Histogram:
//      - req_cost_time
//   Summary:
//      - req_cost_time

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	//ENV      = "env"
	APP = "snow"
	VER = "ver"
)

var (
	collectors = []prometheus.Collector{}
)

func RegisterCollector(c ...prometheus.Collector) {
	collectors = append(collectors, c...)
}

type Options struct {
	labels        map[string]string
	processEnable bool
	runtimeEnable bool
}

type Option func(opt *Options)

// 添加App和Ver label
func AppVer(app, ver string) Option {
	return func(opt *Options) {
		if app != "" {
			opt.labels[APP] = app
		}
		if ver != "" {
			opt.labels[VER] = ver
		}
	}
}

// 添加额外label
func WithLabel(key, val string) Option {
	return func(opt *Options) {
		if key != "" && val != "" {
			opt.labels[key] = val
		}
	}
}

// 收集进程信息
func EnableProcess() Option {
	return func(opt *Options) {
		opt.processEnable = true
	}
}

func EnableRuntime() Option {
	return func(opt *Options) {
		opt.runtimeEnable = true
	}
}

type Reporter struct {
	opts       Options
	collectors []prometheus.Collector
	// registerer
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

var (
	once     sync.Once
	reporter Reporter
)

func Init(opts ...Option) {
	_opts := Options{
		labels: map[string]string{},
	}
	for _, opt := range opts {
		opt(&_opts)
	}

	once.Do(func() {
		cs := collectors
		if _opts.processEnable {
			cs = append(cs, prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		}

		if _opts.runtimeEnable {
			cs = append(cs, prometheus.NewGoCollector())
		}

		reporter = Reporter{
			opts:       _opts,
			collectors: cs,
		}

		registry := prometheus.NewRegistry()

		reporter.registerer = prometheus.WrapRegistererWith(reporter.opts.labels, registry)
		reporter.gatherer = registry

		reporter.registerer.MustRegister(reporter.collectors...)

	})

}

func (p *Reporter) newCounterVec(metric string, labels []string) *prometheus.CounterVec {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metric,
	}, labels)

	return counterVec
}

func (p *Reporter) newGaugeVec(metric string, labels []string) *prometheus.GaugeVec {
	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: metric,
	}, labels)
	return gaugeVec
}

func (p *Reporter) newHistogramVec(metric string, labels []string, buckets []float64) *prometheus.HistogramVec {
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    metric,
		Buckets: buckets,
	}, labels)
	return histogramVec
}

func Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		reporter.registerer, promhttp.HandlerFor(reporter.gatherer, promhttp.HandlerOpts{}),
	)
}
