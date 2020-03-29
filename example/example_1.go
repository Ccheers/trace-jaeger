package main

import (
	"context"
	"fmt"
	"github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/time"
	"github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"kratos/trace_jaeger"
	"os"
	"os/signal"
	"syscall"
	time2 "time"
)

func main() {
	engine := blademaster.DefaultServer(&blademaster.ServerConfig{
		Addr:    "0.0.0.0:8888",
		Timeout: time.Duration(time2.Second),
	})
	engine.Use(trace_jaeger.StartTracer(&jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const", //固定采样
			Param: 1,       //1=全采样、0=不采样
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: true,
			//LocalAgentHostPort: "172.18.0.5:6831",
			LocalAgentHostPort: "127.0.0.1:6831",
		},
		ServiceName: "example_1",
	}))
	engine.GET("/", func(ctx *blademaster.Context) {
		//1.创建子span
		parentSpanContext, _ := ctx.Get("ParentSpanContext")
		trace_jaeger.PushPoint(parentSpanContext.(context.Context), "index_test", "test1", "test1", func() {
			time2.Sleep(time2.Second)
		})
		k := map[string]string{
			"Hello": foo("req2", parentSpanContext.(context.Context)),
		}
		ctx.JSON(k, nil)
	})
	err := engine.Start()
	if err != nil {
		panic(err.Error())
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		fmt.Printf("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Printf("index exit")
			time2.Sleep(time2.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func foo(req string, ctx context.Context) (reply string) {
	//1.创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "foo")
	defer func() {
		//4.接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	//2.模拟处理耗时
	time2.Sleep(time2.Second / 2)
	//3.返回reply
	reply = "foo3Reply"
	return
}
