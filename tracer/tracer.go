package tracer

import (
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerprometheus "github.com/uber/jaeger-lib/metrics/prometheus"

	"service_template/logger"
)

type Config struct {
	AgentAddress string
	Sampler      *SamplerConfig
}

type SamplerConfig struct {
	Type  string
	Param float64
}

func Init(serviceName string, config Config, l logger.Logger) {
	cfg := jaegercfg.Configuration{
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: config.AgentAddress,
		},
	}

	if config.Sampler != nil {
		cfg.Sampler = &jaegercfg.SamplerConfig{
			Type:  config.Sampler.Type,
			Param: config.Sampler.Param,
		}
	}

	options := []jaegercfg.Option{
		jaegercfg.Metrics(jaegerprometheus.New()),
		jaegercfg.Logger(NewJaegerLogger(l)),
	}

	if _, err := cfg.InitGlobalTracer(serviceName, options...); err != nil {
		panic(err)
	}
}

type JaegerLogger struct {
	l logger.Logger
}

func NewJaegerLogger(l logger.Logger) *JaegerLogger {
	return &JaegerLogger{l: l.WithField("from", "jaeger")}
}

var _ jaeger.Logger = &JaegerLogger{}

func (l *JaegerLogger) Error(msg string) {
	l.l.Errorf(msg)
}

func (l *JaegerLogger) Infof(msg string, args ...interface{}) {
	l.l.Infof(msg, args...)
}
