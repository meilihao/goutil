package tracer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/micro/go-micro/metadata"
	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type JaegerConfig struct {
	Name      string
	AgentAddr string
	Enable    bool
	LogSpans  bool
}

func InitTracer(conf *JaegerConfig, logger *zap.SugaredLogger) (io.Closer, error) {
	if !conf.Enable {
		opentracing.SetGlobalTracer(opentracing.NoopTracer{})
		return ioutil.NopCloser(nil), nil
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            conf.LogSpans,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  conf.AgentAddr,
		},
	}

	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.New(
		conf.Name,
		jaegercfg.Logger(jaegerLoggerAdapter{logger}),
		jaegercfg.Metrics(jMetricsFactory),
	)

	if err != nil {
		return ioutil.NopCloser(nil), err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

type jaegerLoggerAdapter struct {
	logger *zap.SugaredLogger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Errorw(msg, zap.String("service", "jaeger"))
}
func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Infow(fmt.Sprintf(msg, args...), zap.String("service", "jaeger"))
}

func GenerateSpanByHeader(ctx context.Context, service, method string, req interface{}) (opentracing.Span, metadata.Metadata) {
	var sp opentracing.Span

	md, ok := metadata.FromContext(ctx)
	if !ok {
		sp = opentracing.StartSpan(service + "#" + method)

		return sp, nil
	}

	traceId := md["Uber-Trace-Id"]
	if traceId == "" {
		sp = opentracing.StartSpan(service + "#" + method)

		return sp, md
	}

	pSp, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier{"Uber-Trace-Id": traceId})
	if err != nil {
		sp = opentracing.StartSpan(service + "#" + method)

		sp.SetTag("error", true)
		sp.LogKV("error", err.Error())
	} else {
		sp = opentracing.StartSpan(service+"#"+method, opentracing.ChildOf(pSp))
	}

	return sp, md
}

func getApiName(c echo.Context) string {
	return fmt.Sprintf("%s#%s", c.Request().RequestURI, c.Request().Method)
}

func GenerateSpan(c echo.Context) (opentracing.Span, context.Context) {
	sp := opentracing.StartSpan(getApiName(c))

	r := c.Request()
	sp.SetTag("http.url", r.URL.String())

	if r.URL.Path != "/v1/user/login" { // 密码脱敏
		var reqBody []byte
		reqBody, _ = ioutil.ReadAll(r.Body)
		r.Body.Close()

		if len(reqBody) > 0 {
			bf := bytes.NewBuffer(reqBody)
			r.Body = ioutil.NopCloser(bf)

			sp.SetTag("http.body", string(reqBody))
		}
	}

	h := opentracing.TextMapCarrier{}
	opentracing.GlobalTracer().Inject(sp.Context(), opentracing.TextMap, h)

	// echo request id
	tracdId := ExtractUberTraceId(h["uber-trace-id"])
	c.Response().Header().Set(echo.HeaderXRequestID, tracdId)
	h["id"] = tracdId

	mdc := metadata.NewContext(context.Background(), metadata.Metadata(h))

	return sp, mdc
}

func ExtractUberTraceId(id string) string {
	if id == "" {
		return ""
	}

	i := strings.Index(id, ":")

	return id[:i]
}
