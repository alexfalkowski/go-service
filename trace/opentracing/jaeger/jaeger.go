package jaeger

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	ozap "github.com/alexfalkowski/go-service/trace/opentracing/logger/zap"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	jprometheus "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	eventsPerSecond = 100
)

// NewTracer for jaeger.
// nolint:ireturn
func NewTracer(lc fx.Lifecycle, name string, logger *zap.Logger, cfg *Config) (opentracing.Tracer, error) {
	c := jconfig.Configuration{
		ServiceName: name,
		Sampler: &jconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: eventsPerSecond,
		},
		Reporter: &jconfig.ReporterConfig{
			LocalAgentHostPort: cfg.Host,
			LogSpans:           false,
		},
	}

	options := []jconfig.Option{
		jconfig.Logger(ozap.NewLogger(logger)),
		jconfig.Metrics(jprometheus.New()),
	}

	tracer, closer, err := c.NewTracer(options...)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return closer.Close()
		},
	})

	return tracer, nil
}

// Register for jaeger.
func Register(lc fx.Lifecycle, logger *zap.Logger, cfg *Config) error {
	tracer, err := NewTracer(lc, os.ExecutableName(), logger, cfg)
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	return nil
}
