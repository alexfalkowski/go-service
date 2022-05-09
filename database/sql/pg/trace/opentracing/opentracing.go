package opentracing

import (
	"context"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/version"
	otr "github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

// TracerParams for opentracing.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for opentracing.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "pg", Config: params.Config, Version: params.Version})
}

// Tracer for opentracing.
type Tracer otr.Tracer

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "pg", operation, method, opts...)
}
