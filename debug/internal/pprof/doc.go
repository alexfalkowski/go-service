// Package pprof provides Fx wiring to expose net/http/pprof profiling endpoints.
//
// This package integrates pprof handlers into the debug HTTP server setup so you can
// capture CPU, heap, goroutine, mutex, and block profiles for diagnostics.
//
// It is registered by the root debug package.
package pprof
