package telemetry_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/cache/telemetry"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestInstrumentTracingDisablesRawCommandStatements(t *testing.T) {
	test.ResetTelemetry(t)
	t.Cleanup(func() {
		test.ResetTelemetry(t)
	})

	exporter := &test.SpanExporter{}
	provider := tracer.NewProvider(tracer.WithSyncer(exporter))
	tracer.SetProvider(provider)
	t.Cleanup(func() {
		require.NoError(t, provider.Shutdown(context.Background()))
	})

	client := redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:1",
		DialTimeout:  (10 * time.Millisecond).Duration(),
		ReadTimeout:  (10 * time.Millisecond).Duration(),
		MaxRetries:   0,
		WriteTimeout: (10 * time.Millisecond).Duration(),
	})
	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})

	require.NoError(t, telemetry.InstrumentTracing(client))

	const secretKey = "secret-cache-key"
	_ = client.Get(t.Context(), secretKey).Err()

	spans := exporter.Spans()
	require.NotEmpty(t, spans)

	for _, span := range spans {
		for _, attr := range span.Attributes() {
			require.NotEqual(t, "db.statement", string(attr.Key))
			require.NotEqual(t, "db.query.text", string(attr.Key))
			require.False(t, strings.Contains(attr.Value.AsString(), secretKey), "redis trace attribute leaked command key")
		}
	}
}
