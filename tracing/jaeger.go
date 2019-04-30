package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

// NewJaegerTracer ...
func NewJaegerTracer(serviceName, url string, log jaeger.Logger) opentracing.Tracer {
	sender := transport.NewHTTPTransport(
		url,
	)
	tracer, _ := jaeger.NewTracer(serviceName,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(log)),
	)
	return tracer
}
