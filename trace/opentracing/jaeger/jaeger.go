package jaeger

import (
	"context"

	ozap "github.com/alexfalkowski/go-service/trace/opentracing/logger/zap"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
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
		jconfig.Metrics(jprometheus.New(jprometheus.WithRegisterer(prometheus.NewRegistry()))),
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
