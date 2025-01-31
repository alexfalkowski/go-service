package tracer

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Meta for tracer.
func Meta(ctx context.Context, span trace.Span) {
	strs := meta.SnakeStrings(ctx, "meta.")
	attrs := make([]attribute.KeyValue, len(strs))
	cnt := 0

	for k, v := range strs {
		attrs[cnt] = attribute.Key(k).String(v)
		cnt++
	}

	span.SetAttributes(attrs...)
}

// WithTraceID for tracer.
func WithTraceID(ctx context.Context, span trace.Span) context.Context {
	return meta.WithAttribute(ctx, "traceId", meta.ToString(span.SpanContext().TraceID()))
}
