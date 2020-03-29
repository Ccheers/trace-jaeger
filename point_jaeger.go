package trace_jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	openTracingLog "github.com/opentracing/opentracing-go/log"
)

func PushPoint(ctx context.Context, operationName, event, val string, f func(), opts ...opentracing.StartSpanOption) {
	span, _ := opentracing.StartSpanFromContext(ctx, operationName, opts...)
	span.LogFields(
		openTracingLog.String("event", event),
		openTracingLog.String("value", val),
	)
	f()
	span.Finish()
}
