package main

const (
	_tplUtil = ``

	_tplUtilsMetric = `package metric

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
`
	_tplUtilsHttpClient = `package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"{{.ModuleName}}/app/http/trace"
	"{{.ModuleName}}/config"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"github.com/go-resty/resty/v2"
	"github.com/qit-team/snow-core/log/logger"
)

const (
	RetryCounts   = 2
	RetryInterval = 3 * time.Second
)

const componentIDGOHttpClient = 5005

type ClientConfig struct {
	name      string
	ctx       context.Context
	client    *resty.Client
	tracer    *go2sky.Tracer
	extraTags map[string]string
}

type ClientOption func(*ClientConfig)

func WithClientTag(key string, value string) ClientOption {
	return func(c *ClientConfig) {
		if c.extraTags == nil {
			c.extraTags = make(map[string]string)
		}
		c.extraTags[key] = value
	}
}

func WithClient(client *resty.Client) ClientOption {
	return func(c *ClientConfig) {
		c.client = client
	}
}

func WithContext(ctx context.Context) ClientOption {
	return func(c *ClientConfig) {
		c.ctx = ctx
	}
}

type transport struct {
	*ClientConfig
	delegated http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	span, err := t.tracer.CreateExitSpan(t.ctx, fmt.Sprintf("/%s%s", req.Method, req.URL.Path), req.Host, func(header string) error {
		req.Header.Set(propagation.Header, header)
		return nil
	})
	if err != nil {
		return t.delegated.RoundTrip(req)
	}
	defer span.End()
	span.SetComponent(componentIDGOHttpClient)
	for k, v := range t.extraTags {
		span.Tag(go2sky.Tag(k), v)
	}
	span.Tag(go2sky.TagHTTPMethod, req.Method)
	span.Tag(go2sky.TagURL, req.URL.String())
	span.SetSpanLayer(v3.SpanLayer_Http)
	resp, err = t.delegated.RoundTrip(req)
	if err != nil {
		span.Error(time.Now(), err.Error())
		return
	}
	span.Tag(go2sky.TagStatusCode, strconv.Itoa(resp.StatusCode))
	if resp.StatusCode >= http.StatusBadRequest {
		span.Error(time.Now(), "Errors on handling client")
	}
	return resp, nil
}

func NewClient(ctx context.Context, options ...ClientOption) (client *resty.Client) {
	client = resty.New()
	if config.IsDebug() {
		client.SetDebug(true).EnableTrace()
	}

	var (
		tracer *go2sky.Tracer
		err    error
	)
	if len(config.GetConf().SkyWalkingOapServer) > 0 && config.IsEnvEqual(config.ProdEnv) {
		tracer, err = trace.Tracer()
		if err != nil {
			logger.Error(ctx, "NewClient:Tracer", err.Error())
		}
	}
	if tracer != nil {
		co := &ClientConfig{ctx: ctx, tracer: tracer}
		for _, option := range options {
			option(co)
		}
		if co.client == nil {
			co.client = client
		}
		tp := &transport{
			ClientConfig: co,
			delegated:    http.DefaultTransport,
		}
		if co.client.GetClient().Transport != nil {
			tp.delegated = co.client.GetClient().Transport
		}
		co.client.SetTransport(tp)
	}

	client.OnBeforeRequest(func(ct *resty.Client, req *resty.Request) error {
		//req.SetContext(c)
		logger.Info(ctx, "OnBeforeRequest", logger.NewWithField("url", req.URL))
		return nil // if its success otherwise return error
	})
	// Registering Response Middleware
	client.OnAfterResponse(func(ct *resty.Client, resp *resty.Response) error {
		logger.Info(ctx, "OnAfterResponse", logger.NewWithField("url", resp.Request.URL), logger.NewWithField("request", resp.Request.RawRequest), logger.NewWithField("response", resp.String()))
		return nil
	})
	return client
}

func NewClientWithRetry(ctx context.Context, retryCounts int, retryInterval time.Duration, options ...ClientOption) (client *resty.Client) {
	client = resty.New()
	if config.IsDebug() {
		client.SetDebug(true).EnableTrace()
	}
	if retryCounts == 0 {
		retryCounts = RetryCounts
	}
	if retryInterval.Seconds() == 0.0 {
		retryInterval = RetryInterval
	}
	client.SetRetryCount(retryCounts).SetRetryMaxWaitTime(retryInterval)

	var (
		tracer *go2sky.Tracer
		err    error
	)
	if len(config.GetConf().SkyWalkingOapServer) > 0 && config.IsEnvEqual(config.ProdEnv) {
		tracer, err = trace.Tracer()
		if err != nil {
			logger.Error(ctx, "NewClient:Tracer", err.Error())
		}
	}
	if tracer != nil {
		co := &ClientConfig{ctx: ctx, tracer: tracer}
		for _, option := range options {
			option(co)
		}
		if co.client == nil {
			co.client = client
		}
		tp := &transport{
			ClientConfig: co,
			delegated:    http.DefaultTransport,
		}
		if co.client.GetClient().Transport != nil {
			tp.delegated = co.client.GetClient().Transport
		}
		co.client.SetTransport(tp)
	}

	client.OnBeforeRequest(func(ct *resty.Client, req *resty.Request) error {
		logger.Info(ctx, "OnBeforeRequest", logger.NewWithField("url", req.URL))
		return nil // if its success otherwise return error
	})
	// Registering Response Middleware
	client.OnAfterResponse(func(ct *resty.Client, resp *resty.Response) error {
		logger.Info(ctx, "OnAfterResponse", logger.NewWithField("url", resp.Request.URL), logger.NewWithField("request", resp.Request.RawRequest), logger.NewWithField("response", resp.String()))
		return nil
	})
	return client
}
`
)
