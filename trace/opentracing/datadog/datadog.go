package datadog

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/trace/opentracing/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// NewTracer for datadog.
// nolint:ireturn
func NewTracer(lc fx.Lifecycle, cfg *Config) opentracing.Tracer {
	opts := []tracer.StartOption{
		tracer.WithService(os.ExecutableName()),
		tracer.WithAgentAddr(cfg.Host),
		tracer.WithLogger(logger.NewLogger()),
	}
	t := opentracer.New(opts...)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			tracer.Stop() // important for data integrity (flushes any leftovers)

			return nil
		},
	})

	return t
}
