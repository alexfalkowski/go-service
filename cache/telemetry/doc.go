// Package telemetry exposes selected Redis OpenTelemetry helpers through the
// go-service cache import tree.
//
// This package is the cache-side wrapper boundary for Redis OpenTelemetry
// instrumentation. It keeps cache code on a go-service import path while
// preserving the behavior of the upstream redisotel helpers used to instrument
// Redis clients for tracing and metrics.
//
// Use this package when instrumenting go-redis clients that back the go-service
// cache subsystem. Higher-level cache code should generally prefer
// `cache/driver.NewDriver`, which applies this instrumentation automatically for
// the built-in Redis backend.
package telemetry
