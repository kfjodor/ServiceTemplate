package logger

import (
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

func WithTraceId(l Logger, span opentracing.Span) Logger {
	sc, ok := span.Context().(jaeger.SpanContext)
	if !ok {
		return l
	}

	return l.WithField("traceId", sc.TraceID().String())
}
