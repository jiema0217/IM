package middleware

import (
	"IMProject/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
)

func Jaeger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader("im-trace-id")
		var span opentracing.Span
		if traceId != "" {
			var err error
			span, err = getParentSpan(ctx.Request.URL.Path, traceId, ctx.Request.Header)
			if err != nil {
				return
			}
		} else {
			span = opentracing.GlobalTracer().StartSpan(ctx.Request.URL.Path)
		}
		defer span.Finish()
		ctx.Set(traceId, opentracing.ContextWithSpan(ctx, span))
		ctx.Next()
	}
}

func InitJaeger() (opentracing.Tracer, io.Closer) {
	cfg := &jaegerCfg.Configuration{
		ServiceName: config.Cfg.ServiceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer(jaegerCfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func getParentSpan(spanName string, traceId string, header http.Header) (opentracing.Span, error) {
	carrier := opentracing.HTTPHeadersCarrier{}
	carrier.Set("im-trace-id", traceId)

	tracer := opentracing.GlobalTracer()
	wireContext, err := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header),
	)

	parentSpan := opentracing.StartSpan(
		spanName,
		ext.RPCServerOption(wireContext))
	if err != nil {
		return nil, err
	}
	return parentSpan, err
}
