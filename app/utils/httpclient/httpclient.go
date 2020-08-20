package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"github.com/go-resty/resty/v2"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow/app/http/trace"
	"github.com/qit-team/snow/config"
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
