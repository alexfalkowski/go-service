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
//
// This mutates the process-global OpenTelemetry tracer provider through
// telemetry/tracer registration. The supplied lifecycle owns provider shutdown.
//
// It panics through [runtime.Must] if tracer registration fails.
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
//
// This mutates the process-global OpenTelemetry tracer provider for the current
// test and resets it to a no-op provider with tb.Cleanup.
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
//
// It resets process-global telemetry before installation and again with
// tb.Cleanup, making it the preferred helper for tests that assert exported
// spans.
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
