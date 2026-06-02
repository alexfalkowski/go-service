// Package telemetry exposes selected Redis OpenTelemetry helpers through the
// go-service cache import tree.
//
// This package is the cache-side wrapper boundary for Redis OpenTelemetry
// instrumentation. It keeps cache code on a go-service import path while
// preserving the behavior of the upstream [github.com/redis/go-redis/extra/redisotel/v9] helpers used to instrument
// Redis clients for tracing and metrics.
//
// Use this package when instrumenting [github.com/redis/go-redis/v9] clients that back the go-service
// cache subsystem. Higher-level cache code should generally prefer
// [github.com/alexfalkowski/go-service/v2/cache/driver.NewDriver], which applies this instrumentation automatically for
// the built-in Redis backend.
package telemetry
