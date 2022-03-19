package jaeger

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	jaegerPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/fx"
)

const (
	eventsPerSecond = 100
)

// Register for jaeger.
func Register(lc fx.Lifecycle, cfg *Config) error {
	name, err := os.ExecutableName()
	if err != nil {
		return err
	}

	c := jaegerConfig.Configuration{
		ServiceName: name,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: eventsPerSecond,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: cfg.Host,
			LogSpans:           false,
		},
	}

	options := []jaegerConfig.Option{
		jaegerConfig.Logger(jaegerLog.NullLogger),
		jaegerConfig.Metrics(jaegerPrometheus.New()),
	}

	tracer, closer, err := c.NewTracer(options...)
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return closer.Close()
		},
	})

	return nil
}
