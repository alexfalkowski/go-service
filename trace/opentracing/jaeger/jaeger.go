package jaeger

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	jlog "github.com/uber/jaeger-client-go/log"
	jprometheus "github.com/uber/jaeger-lib/metrics/prometheus"
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
		jconfig.Logger(jlog.NullLogger),
		jconfig.Metrics(jprometheus.New()),
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
