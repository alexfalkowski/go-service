package datadog

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/trace/opentracing/logger"
	"github.com/alexfalkowski/go-service/version"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TracerParams for datadog.
type TracerParams struct {
	Lifecycle fx.Lifecycle
	Name      string
	Version   version.Version
	Host      string
}

// NewTracer for datadog.
func NewTracer(params TracerParams) opentracing.Tracer {
	opts := []tracer.StartOption{
		tracer.WithService(params.Name),
		tracer.WithAgentAddr(params.Host),
		tracer.WithLogger(logger.NewLogger()),
		tracer.WithGlobalTag("name", os.ExecutableName()),
		tracer.WithGlobalTag("version", params.Version),
	}
	t := opentracer.New(opts...)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			tracer.Stop() // important for data integrity (flushes any leftovers)

			return nil
		},
	})

	return t
}
