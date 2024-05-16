package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Meta for tracer.
func Meta(ctx context.Context, span trace.Span) {
	for k, v := range meta.SnakeStrings(ctx, "meta.") {
		span.SetAttributes(attribute.Key(k).String(v))
	}
}

// WithTraceID for tracer.
func WithTraceID(ctx context.Context, span trace.Span) context.Context {
	if !span.IsRecording() {
		return ctx
	}

	return meta.WithAttribute(ctx, "traceId", meta.ToString(span.SpanContext().TraceID()))
}
