package jaeger

import (
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"go.uber.org/fx"
)

// NewTracer for jaeger.
func NewTracer(lc fx.Lifecycle, cfg *jaeger.Config) (opentracing.Tracer, error) {
	return jaeger.NewTracer(jaeger.TracerParams{Lifecycle: lc, Name: "nsq", Config: cfg})
}
