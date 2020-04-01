package trace_jaeger

import (
	"github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

func StartTracer(conf *jaegerConfig.Configuration) blademaster.HandlerFunc {
	return func(c *blademaster.Context) {
		var parentSpan opentracing.Span

		tracer, closer, err := NewTracer(conf)
		if err != nil {
			panic(err)
		}
		defer func() {
			closer.Close()
		}()

		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}
		//c.Context = context.WithValue(c,"ParentSpanContext",parentSpan.Context())
		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", opentracing.ContextWithSpan(c, parentSpan))
		c.Set("SpanContext", parentSpan.Context())
		c.Next()
	}
}
