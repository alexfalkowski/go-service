package datadog

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Register for datadog.
func Register(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) error {
	opts := []tracer.StartOption{
		tracer.WithService(cfg.AppName),
		tracer.WithAgentAddr(cfg.DataDogTraceHost),
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
