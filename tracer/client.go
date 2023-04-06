package tracer

import (
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/twitchtv/twirp"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TraceHTTPClient struct {
	client HTTPClient
	tracer opentracing.Tracer
}

var _ HTTPClient = (*TraceHTTPClient)(nil)

func (c *TraceHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	methodName, ok := twirp.MethodName(ctx)
	if !ok {
		methodName = req.URL.Path
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, methodName, ext.SpanKindRPCClient)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())

	err := c.tracer.Inject(span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		span.LogFields(otlog.String("event", "tracer.Inject() failed"), otlog.Error(err))
	}
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		setErrorSpan(span, err.Error())
		span.Finish()
		return res, err
	}
	ext.HTTPStatusCode.Set(span, uint16(res.StatusCode))

	if res.StatusCode >= 400 {
		span.SetTag("error", true)
	}

	res.Body = closer{
		ReadCloser: res.Body,
		span:       span,
	}
	return res, nil
}

type closer struct {
	io.ReadCloser
	span opentracing.Span
}

func (c closer) Close() error {
	err := c.ReadCloser.Close()
	c.span.Finish()
	return err
}

func setErrorSpan(span opentracing.Span, errorMessage string) {
	span.SetTag("error", true)
	span.LogFields(otlog.String("event", "error"), otlog.String("message", errorMessage))
}
