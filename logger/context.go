package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sdkLog "github.com/nilorg/sdk/log"
	"github.com/nilorg/sdk/log/trace"
	"google.golang.org/grpc/metadata"
)

const (
	// XTraceIDKey ...
	XTraceIDKey = "X-Trace-Id"
	// XSpanIDKey ...
	XSpanIDKey = "X-Span-Id"
)

// WithGinContext ...
func WithGinContext(ctx *gin.Context) context.Context {
	parent := ctx.Request.Context()
	if traceID := ctx.GetHeader(XTraceIDKey); traceID != "" {
		parent = sdkLog.NewTraceIDContext(parent, traceID)
	} else {
		parent = sdkLog.NewTraceIDContext(parent, uuid.New().String())
	}
	if spanID := ctx.GetHeader(XSpanIDKey); spanID != "" {
		spanID = trace.StartSpanID(spanID)
		parent = sdkLog.NewSpanIDContext(parent, spanID)
	} else {
		parent = sdkLog.NewSpanIDContext(parent, "0")
	}
	return parent
}

// WithGrpcMetadata 从上下文中
func WithGrpcMetadata(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	if v := md.Get(TraceIDKey); len(v) > 0 {
		ctx = sdkLog.NewTraceIDContext(ctx, v[0])
	}
	if v := md.Get(SpanIDKey); len(v) > 0 {
		spanID := trace.StartSpanID(v[0])
		ctx = sdkLog.NewSpanIDContext(ctx, spanID)
	}
	return ctx
}
