package infra

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"io"
	"net"
	"net/http"
	"time"

	"service_template/logger"
	"service_template/tracer"
)

type Config struct {
	ServiceName             string `validate:"nonzero"`
	Logger                  *logger.Config
	LoggerHooks             []logrus.Hook
	Tracer                  *tracer.Config
	GracefulShutdownTimeout time.Duration
}

const defaultGracefulShutdownTimeout = 10 * time.Second

type healthContextKey struct{}
type configContextKey struct{}
type traceContextKey struct{}

var ErrInvalidContext = errors.New("invalid context")

// GenerateNewTrace
func GenerateNewTrace(ctx context.Context) context.Context {
	var serverSpan opentracing.Span
	serverSpan = opentracing.StartSpan("HTTP request")
	defer serverSpan.Finish()

	if sc, ok := serverSpan.Context().(jaeger.SpanContext); ok {
		traceID := sc.TraceID().String()
		ctx = context.WithValue(ctx, traceContextKey{}, traceID)
	}
	log := logger.FromContext(ctx)
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	log = logger.WithTraceId(log, serverSpan)
	return logger.ToContext(ctx, log)
}

// Context
func Context(config Config, wr io.Writer) context.Context {
	ctx, _ := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, configContextKey{}, config)

	if config.Logger != nil {
		ctx = logger.ToContext(ctx, logger.New(*config.Logger, wr, config.LoggerHooks...))
	}

	log := logger.FromContext(ctx)

	if config.Tracer != nil {
		tracer.Init(config.ServiceName, *config.Tracer, logger.FromContext(ctx))
	}

	//sigCh := make(chan os.Signal, 1)
	//signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	healthStatus := int32(http.StatusTeapot)
	ctx = context.WithValue(ctx, healthContextKey{}, &healthStatus)

	var serverSpan opentracing.Span
	serverSpan = opentracing.StartSpan("HTTP request")
	defer serverSpan.Finish()

	if sc, ok := serverSpan.Context().(jaeger.SpanContext); ok {
		traceID := sc.TraceID().String()
		ctx = context.WithValue(ctx, traceContextKey{}, traceID)
	}

	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	log = logger.WithTraceId(log, serverSpan)
	ctx = logger.ToContext(ctx, log)
	return ctx
}

// ServeHTTP
func ServeHTTP(ctx context.Context, listen string, handler http.Handler) error {
	config, ok := ctx.Value(configContextKey{}).(Config)
	if !ok {
		return ErrInvalidContext
	}

	log := logger.FromContext(ctx).WithField("m", "ServeHTTP")
	log.Debugf("ServeHTTP:: listen: %v, handler: %v", listen, handler)

	var listener net.Listener
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	log.Debugf("listening HTTP requests on %s", listen)

	server := http.Server{Handler: handler}

	go func() {
		log.Debugf("HTTP server exit message: %v", server.Serve(listener))
	}()

	<-ctx.Done()

	shutdownTimeout := config.GracefulShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = defaultGracefulShutdownTimeout
	}

	shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelFunc()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Debugf("HTTP server graceful shutdown message: %v", err)
	}

	return nil
}
