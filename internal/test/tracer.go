package test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
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

// EnableSpanExporter installs a synchronous tracer provider backed by a SpanExporter.
func EnableSpanExporter(tb testing.TB) *SpanExporter {
	tb.Helper()

	exporter := &SpanExporter{}
	provider := tracer.NewProvider(tracer.WithSyncer(exporter))
	tracer.SetProvider(provider)

	tb.Cleanup(func() {
		require.NoError(tb, provider.Shutdown(context.Background()))
		tracer.SetProvider(tracer.NewNoopProvider())
	})

	return exporter
}

// EnableIsolatedSpanExporter resets telemetry and installs a synchronous span exporter.
func EnableIsolatedSpanExporter(tb testing.TB) *SpanExporter {
	tb.Helper()

	ResetTelemetry(tb)
	tb.Cleanup(func() {
		ResetTelemetry(tb)
	})

	return EnableSpanExporter(tb)
}

// SpanExporter records exported spans for test assertions.
type SpanExporter struct {
	spans []tracer.ReadOnlySpan
	mu    sync.Mutex
}

// ExportSpans records the supplied spans.
func (e *SpanExporter) ExportSpans(_ context.Context, spans []tracer.ReadOnlySpan) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.spans = append(e.spans, spans...)

	return nil
}

// Shutdown satisfies [tracer.SpanExporter].
func (e *SpanExporter) Shutdown(context.Context) error {
	return nil
}

// Spans returns a copy of the recorded spans.
func (e *SpanExporter) Spans() []tracer.ReadOnlySpan {
	e.mu.Lock()
	defer e.mu.Unlock()

	return append([]tracer.ReadOnlySpan(nil), e.spans...)
}
