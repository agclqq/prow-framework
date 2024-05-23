package skywalking

import (
	"context"
	"fmt"

	"github.com/SkyAPM/go2sky"
	"github.com/gin-gonic/gin"

	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

var tracer *go2sky.Tracer

func GetTraceId(ctx context.Context) string {
	if c, ok := ctx.(*gin.Context); ok {
		traceId := go2sky.TraceID(c.Request.Context())
		return traceId
	}
	traceId := go2sky.TraceID(ctx)
	return traceId
}
func GetFullTraceInfo(ctx context.Context) string {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	serviceName := go2sky.ServiceName(ctx)
	serviceInstanceName := go2sky.ServiceInstanceName(ctx)
	traceId := go2sky.TraceID(ctx)
	traceSegmentID := go2sky.TraceSegmentID(ctx)
	spanID := go2sky.SpanID(ctx)

	return fmt.Sprintf("serviceName:%s-serviceInstanceName:%s-traceId:%s-traceSegmentID:%s-spanID:%d", serviceName, serviceInstanceName, traceId, traceSegmentID, spanID)
}
func SetTracer(trac *go2sky.Tracer) {
	tracer = trac
}
func GetTracer(opts ...go2sky.TracerOption) *go2sky.Tracer {
	//var err error
	//if tracer == nil {
	//	tracer, err = go2sky.NewTracer("publish-server1", opts...)
	//}
	//if err != nil {
	//	return nil
	//}
	return tracer
}

func AddTrace(ctx context.Context, tracer *go2sky.Tracer, operationName string) (go2sky.Span, context.Context, error) {
	localSpan, ctx, err := tracer.CreateEntrySpan(ctx, operationName, func(headerKey string) (string, error) {
		if v, ok := ctx.Value("trace-" + headerKey).(string); ok {
			return v, nil
		}
		return "", nil
	})
	if err != nil {
		return nil, ctx, err
	}
	localSpan.SetSpanLayer(agentv3.SpanLayer_Unknown)
	return localSpan, ctx, nil
}
