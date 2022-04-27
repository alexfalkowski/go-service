package datadog

import (
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/transport/http/trace/opentracing"
	"go.uber.org/fx"
)

// NewTracer for datadog.
func NewTracer(lc fx.Lifecycle, cfg *datadog.Config) opentracing.Tracer {
	return datadog.NewTracer(datadog.TracerParams{Lifecycle: lc, Name: "http", Config: cfg})
}
