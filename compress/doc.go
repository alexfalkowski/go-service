// Package compress provides compression abstractions and wiring for go-service.
//
// This package defines a common Compressor interface and provides Fx wiring that constructs a map of
// supported compressors (zstd, s2, snappy, and none) keyed by kind.
//
// Start with `Compressor`, `Map`, and `Module`.
package compress
