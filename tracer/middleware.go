package tracer

import (
	"bytes"
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	ext "github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"

	"service_template/logger"
)

type LogResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *LogResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

func tracer(w http.ResponseWriter, r *http.Request) context.Context {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))

	if err != nil {
		log.Debugf("failed to extract span from headers")
	}

	var serverSpan opentracing.Span

	options := []opentracing.StartSpanOption{
		ext.RPCServerOption(wireContext),
	}

	if wireContext != nil {
		options = append(options, opentracing.ChildOf(wireContext))
	}

	serverSpan = opentracing.StartSpan("HTTP request", options...)
	defer serverSpan.Finish()

	if sc, ok := serverSpan.Context().(jaeger.SpanContext); ok {
		headers := w.Header()
		traceID := sc.TraceID().String()
		headers.Set("X-Trace-ID", traceID)
	}

	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	log = logger.WithTraceId(log, serverSpan)
	ctx = logger.ToContext(ctx, log)

	return ctx
}
