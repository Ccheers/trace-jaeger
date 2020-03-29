package trace_jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
)

func NewTracer(Conf *jaegerConfig.Configuration) (opentracing.Tracer, io.Closer, error) {
	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := Conf.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}
