// Package bytes provides byte-oriented helpers used across go-service.
//
// Most identifiers in this package are thin wrappers or aliases around the standard library `bytes`
// package. They exist to:
//   - provide a consistent import path within go-service, and
//   - centralize a small set of byte/string conversion utilities.
//
// # Aliases and wrappers
//
// `Buffer` is an alias for `bytes.Buffer`, constructors like `NewBuffer`, `NewBufferString`, and
// `NewReader` delegate directly to the standard library, and helpers like `Clone` and `TrimSpace`
// mirror standard library behavior through the go-service import path.
//
// # Zero-copy conversions
//
// This package also exposes `String`, which converts a `[]byte` to a `string` without allocating.
// This is an advanced, performance-oriented helper with important safety constraints; see `String`
// for details.
//
// # Human-readable sizes
//
// `Size` is a named byte-count type for typed configuration surfaces. `Size.String` and `ParseSize`
// use human-readable decimal size strings such as `64B`, `2MB`, and `4GB`.
//
// Text and JSON marshaling emit exact raw byte counts with a `B` suffix, such as `4000000B`.
// Text and JSON unmarshaling accept the same decimal size inputs as `ParseSize`.
package bytes
