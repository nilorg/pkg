package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// ContextWithSpan ...
func ContextWithSpan(span opentracing.Span) context.Context {
	ctx := context.Background()
	if span != nil {
		// 设置parent span
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return ctx
}
