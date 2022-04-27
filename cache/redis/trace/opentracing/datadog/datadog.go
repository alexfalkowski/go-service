package datadog

import (
	"github.com/alexfalkowski/go-service/cache/redis/trace/opentracing"
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"go.uber.org/fx"
)

// NewTracer for datadog.
func NewTracer(lc fx.Lifecycle, cfg *datadog.Config) opentracing.Tracer {
	return datadog.NewTracer(datadog.TracerParams{Lifecycle: lc, Name: "redis", Config: cfg})
}
