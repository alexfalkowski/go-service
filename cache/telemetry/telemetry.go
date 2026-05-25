package telemetry

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	client "github.com/redis/go-redis/v9"
)

// InstrumentTracing instruments a Redis client for OpenTelemetry tracing.
//
// This is a thin wrapper around redisotel.InstrumentTracing with raw command
// statement capture disabled.
//
// The provided client is modified in place to emit tracing data for supported
// Redis operations. The wrapper does not change upstream behavior or error
// semantics.
func InstrumentTracing(client client.UniversalClient) error {
	return redisotel.InstrumentTracing(client, redisotel.WithDBStatement(false))
}

// InstrumentMetrics instruments a Redis client for OpenTelemetry metrics.
//
// This is a thin wrapper around redisotel.InstrumentMetrics that unregisters
// observable callbacks when closeChan is closed.
//
// The provided client is modified in place to emit Redis client metrics. The
// wrapper does not change upstream behavior or error semantics.
func InstrumentMetrics(client client.UniversalClient, closeChan chan struct{}) error {
	return redisotel.InstrumentMetrics(client, redisotel.WithCloseChan(closeChan))
}
