package zlog

import (
	"context"
	"net/http"

	"github.com/nilorg/pkg/zlog/trace"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Standard *zap.Logger
	Sugared  *zap.SugaredLogger
)

type Options struct {
	Config zap.Config
}
type Option func(o *Options)

func newOptions(opts ...Option) Options {
	o := Options{
		Config: zap.NewDevelopmentConfig(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithConfig(conf zap.Config) Option {
	return func(o *Options) {
		o.Config = conf
	}
}

func Init(opts ...Option) {
	trace.Init()
	opt := newOptions(opts...)
	Standard, _ = opt.Config.Build()
	Sugared = Standard.Sugar()
}

func InitForViper(conf *viper.Viper) {
	c := zap.NewDevelopmentConfig()
	if conf.GetString("zlog.mod") == "production" {
		c = zap.NewProductionConfig()
	}
	if conf.GetBool("zlog.zinc.enabled") {
		RegisterSink()
		c.OutputPaths = append(c.OutputPaths, conf.GetString("zlog.zinc.url"))
		c.ErrorOutputPaths = append(c.ErrorOutputPaths, conf.GetString("zlog.zinc.url"))
	}
	Init(WithConfig(c))
}
func RegisterSink() {
	// 将ZincSink工厂函数注册到zap中, 自定义协议名为 zinc
	if err := zap.RegisterSink("zinc", NewZincSink); err != nil {
		panic(err)
	}
}

func Sync() {
	Sugared.Sync()
	Standard.Sync()
}

type (
	ContextTraceIDKey struct{}
	ContextSpanIDKey  struct{}
	ContextUserIDKey  struct{}
)

const (
	// TraceIDKey 跟踪ID
	TraceIDKey = "trace_id"
	// SpanIDKey 请求ID
	SpanIDKey = "span_id"
	// UserIDKey 用户ID
	UserIDKey = "user_id"
)

func With(ctx context.Context) *zap.Logger {
	fields := make([]zap.Field, 0)
	if traceID, ok := FromTraceIDContext(ctx); ok {
		fields = append(fields, zap.String(TraceIDKey, traceID))
	}
	if spanID, ok := FromSpanIDContext(ctx); ok {
		fields = append(fields, zap.String(SpanIDKey, spanID))
	}
	if userID, ok := FromUserIDContext(ctx); ok {
		fields = append(fields, zap.String(UserIDKey, userID))
	}
	if len(fields) > 0 {
		return Standard.With(fields...)
	} else {
		return Standard
	}
}

func WithSugared(ctx context.Context) *zap.SugaredLogger {
	fields := make([]interface{}, 0)
	if traceID, ok := FromTraceIDContext(ctx); ok {
		fields = append(fields, TraceIDKey, traceID)
	}
	if spanID, ok := FromSpanIDContext(ctx); ok {
		fields = append(fields, SpanIDKey, spanID)
	}
	if userID, ok := FromUserIDContext(ctx); ok {
		fields = append(fields, UserIDKey, userID)
	}
	if len(fields) > 0 {
		return Sugared.With(fields...)
	} else {
		return Sugared
	}
}

// NewTraceIDContext ...
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextTraceIDKey{}, traceID)
}

// FromTraceIDContext ...
func FromTraceIDContext(ctx context.Context) (traceID string, ok bool) {
	traceID, ok = ctx.Value(ContextTraceIDKey{}).(string)
	return
}

// NewSpanIDContext ...
func NewSpanIDContext(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, ContextSpanIDKey{}, spanID)
}

// FromSpanIDContext ...
func FromSpanIDContext(ctx context.Context) (spanID string, ok bool) {
	spanID, ok = ctx.Value(ContextSpanIDKey{}).(string)
	return
}

// NewUserIDContext ...
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextUserIDKey{}, userID)
}

// FromUserIDContext ...
func FromUserIDContext(ctx context.Context) (userID string, ok bool) {
	userID, ok = ctx.Value(ContextUserIDKey{}).(string)
	return
}

// CopyContext copy context
func CopyContext(ctx context.Context) context.Context {
	parent := context.Background()
	if traceID, ok := FromTraceIDContext(ctx); ok {
		parent = NewTraceIDContext(parent, traceID)
	}
	if spanID, ok := FromSpanIDContext(ctx); ok {
		parent = NewSpanIDContext(parent, spanID)
	}
	if userID, ok := FromUserIDContext(ctx); ok {
		parent = NewUserIDContext(parent, userID)
	}
	return parent
}

const (
	// XTraceIDKey ...
	XTraceIDKey = "X-Trace-Id"
	// XSpanIDKey ...
	XSpanIDKey = "X-Span-Id"
)

// WithRequestContext ...
func WithRequestContext(req *http.Request, userID string) context.Context {
	parent := req.Context()
	traceID := req.Header.Get(XTraceIDKey)
	if traceID != "" {
		parent = NewTraceIDContext(parent, traceID)
	} else {
		parent = NewTraceIDContext(parent, trace.NewID())
	}
	if spanID := req.Header.Get(XSpanIDKey); spanID != "" {
		spanID = trace.StartSpanID(traceID, spanID)
		parent = NewSpanIDContext(parent, spanID)
	} else {
		parent = NewSpanIDContext(parent, "0")
	}
	if userID != "" {
		parent = NewUserIDContext(parent, userID)
	}
	return parent
}
