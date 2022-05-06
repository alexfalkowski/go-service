package opentracing

import (
	"context"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	otr "github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

// Tracer for opentracing.
type Tracer otr.Tracer

// NewTracer for opentracing.
func NewTracer(lc fx.Lifecycle, cfg *opentracing.Config) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: lc, Name: "pg", Config: cfg})
}

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "pg", operation, method, opts...)
}
