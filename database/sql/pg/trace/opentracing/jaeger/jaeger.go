package jaeger

import (
	"github.com/alexfalkowski/go-service/database/sql/pg/trace/opentracing"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

// NewTracer for jaeger.
func NewTracer(lc fx.Lifecycle, cfg *jaeger.Config) (opentracing.Tracer, error) {
	return jaeger.NewTracer(jaeger.TracerParams{Lifecycle: lc, Name: "pg", Config: cfg})
}
