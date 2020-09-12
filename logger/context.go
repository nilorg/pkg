package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sdkLog "github.com/nilorg/sdk/log"
	"github.com/nilorg/sdk/log/trace"
)

// WithGinContext ...
func WithGinContext(ctx *gin.Context) context.Context {
	parent := ctx.Request.Context()
	if traceID := ctx.GetHeader("X-Trace-Id"); traceID != "" {
		parent = sdkLog.NewTraceIDContext(parent, traceID)
	} else {
		parent = sdkLog.NewTraceIDContext(parent, uuid.New().String())
	}
	if spanID := ctx.GetHeader("X-Span-Id"); spanID != "" {
		spanID = trace.StartSpanID(spanID)
		parent = sdkLog.NewSpanIDContext(parent, spanID)
	} else {
		parent = sdkLog.NewSpanIDContext(parent, "0")
	}
	return parent
}
