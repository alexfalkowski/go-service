package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

// Trace extracts trace correlation attributes from ctx.
//
// When ctx carries a valid OpenTelemetry span context, it returns trace_id and
// span_id string attributes named per the OpenTelemetry log data model so
// stdout log lines can be correlated with their trace. It returns no attributes
// when ctx has no valid span, leaving untraced log records unchanged.
func Trace(ctx context.Context) []slog.Attr {
	span := tracer.SpanContextFromContext(ctx)
	if !span.IsValid() {
		return nil
	}

	return []slog.Attr{
		slog.String(traceIDKey, span.TraceID().String()),
		slog.String(spanIDKey, span.SpanID().String()),
	}
}

// NewTraceHandler wraps handler so that records emitted under a valid span gain
// the trace_id/span_id attributes from [Trace].
//
// The stdout-oriented logger kinds (json/text/tint) use it; the otlp logger
// correlates through the OpenTelemetry logging bridge instead. Wrap the root
// handler so the correlation keys stay at the record root, matching the
// OpenTelemetry log data model (the stdout constructors do not open slog groups).
func NewTraceHandler(handler slog.Handler) slog.Handler {
	return &traceHandler{Handler: handler}
}

type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	if attrs := Trace(ctx); len(attrs) > 0 {
		record.AddAttrs(attrs...)
	}

	return h.Handler.Handle(ctx, record)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name)}
}
