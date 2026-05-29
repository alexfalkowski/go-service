package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-sync"
	"go.opentelemetry.io/otel/sdk/trace"
)

// RegisterTracer installs the shared test tracer provider on the supplied lifecycle.
func RegisterTracer(lc di.Lifecycle, config *tracer.Config) {
	params := tracer.TracerParams{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		Config:      config,
	}

	runtime.Must(tracer.Register(params))
}

// SpanExporter records exported spans for test assertions.
type SpanExporter struct {
	spans []trace.ReadOnlySpan
	mu    sync.Mutex
}

// ExportSpans records the supplied spans.
func (e *SpanExporter) ExportSpans(_ context.Context, spans []trace.ReadOnlySpan) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.spans = append(e.spans, spans...)

	return nil
}

// Shutdown satisfies trace.SpanExporter.
func (e *SpanExporter) Shutdown(context.Context) error {
	return nil
}

// Spans returns a copy of the recorded spans.
func (e *SpanExporter) Spans() []trace.ReadOnlySpan {
	e.mu.Lock()
	defer e.mu.Unlock()

	return append([]trace.ReadOnlySpan(nil), e.spans...)
}
