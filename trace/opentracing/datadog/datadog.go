package datadog

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Register for datadog.
func Register(lc fx.Lifecycle, cfg *Config) error {
	name, err := os.ExecutableName()
	if err != nil {
		return err
	}

	opts := []tracer.StartOption{
		tracer.WithService(name),
		tracer.WithAgentAddr(cfg.Host),
		tracer.WithLogStartup(false),
	}
	t := opentracer.New(opts...)

	opentracing.SetGlobalTracer(t)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			tracer.Stop() // important for data integrity (flushes any leftovers)

			return nil
		},
	})

	return nil
}
