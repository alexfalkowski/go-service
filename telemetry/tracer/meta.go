package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"go.opentelemetry.io/otel/trace"
)

// Meta for tracer.
func Meta(ctx context.Context, span trace.Span) {
	strings := meta.SnakeStrings(ctx, "meta.")

	for k, v := range strings {
		span.SetAttributes(attributes.String(k, v))
	}
}

// WithTraceID for tracer.
func WithTraceID(ctx context.Context, span trace.Span) context.Context {
	return meta.WithAttribute(ctx, "traceId", meta.ToString(span.SpanContext().TraceID()))
}
