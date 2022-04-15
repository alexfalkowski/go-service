package datadog

import (
	"context"

	"github.com/alexfalkowski/go-service/os"
	ozap "github.com/alexfalkowski/go-service/trace/opentracing/logger/zap"
	"github.com/alexfalkowski/go-service/transport/http"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// NewTracer for datadog.
// nolint:ireturn
func NewTracer(lc fx.Lifecycle, name string, logger *zap.Logger, cfg *Config, httpCfg *http.Config) opentracing.Tracer {
	opts := []tracer.StartOption{
		tracer.WithService(name),
		tracer.WithAgentAddr(cfg.Host),
		tracer.WithLogger(ozap.NewLogger(logger)),
		tracer.WithHTTPClient(http.NewClient(httpCfg, logger)),
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

// Register for datadog.
func Register(lc fx.Lifecycle, logger *zap.Logger, cfg *Config, httpCfg *http.Config) {
	opentracing.SetGlobalTracer(NewTracer(lc, os.ExecutableName(), logger, cfg, httpCfg))
}
