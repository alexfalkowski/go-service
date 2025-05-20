package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/meta"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Meta for tracer.
func Meta(ctx context.Context, span trace.Span) {
	strings := meta.SnakeStrings(ctx, "meta.")

	for k, v := range strings {
		span.SetAttributes(attribute.Key(k).String(v))
	}
}

// WithTraceID for tracer.
func WithTraceID(ctx context.Context, span trace.Span) context.Context {
	return meta.WithAttribute(ctx, "traceId", meta.ToString(span.SpanContext().TraceID()))
}
