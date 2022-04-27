package opentracing

import (
	"context"

	"github.com/alexfalkowski/go-service/trace/opentracing"
	otr "github.com/opentracing/opentracing-go"
)

// Tracer for opentracing.
type Tracer otr.Tracer

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer Tracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return opentracing.StartSpanFromContext(ctx, tracer, "pg", operation, method, opts...)
}
