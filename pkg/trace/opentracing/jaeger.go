package opentracing

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
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

// RegisterJaeger for opentracing.
func RegisterJaeger(lc fx.Lifecycle, cfg *config.Config) error {
	c := jaegerConfig.Configuration{
		ServiceName: cfg.AppName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: eventsPerSecond,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: cfg.TraceHost,
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
