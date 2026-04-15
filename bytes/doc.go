// Package bytes provides byte-oriented helpers used across go-service.
//
// Most identifiers in this package are thin wrappers or aliases around the standard library `bytes`
// package. They exist to:
//   - provide a consistent import path within go-service, and
//   - centralize a small set of byte/string conversion utilities.
//
// # Aliases and wrappers
//
// `Buffer` is an alias for `bytes.Buffer`, and constructors like `NewBuffer`, `NewBufferString`,
// and `NewReader` delegate directly to the standard library.
//
// # Zero-copy conversions
//
// This package also exposes `String`, which converts a `[]byte` to a `string` without allocating.
// This is an advanced, performance-oriented helper with important safety constraints; see `String`
// for details.
//
// # Human-readable sizes
//
// `Size` is a named byte-count type that marshals to and from human-readable SI strings such as
// `64B`, `2MB`, and `4GB`. It is intended for typed configuration surfaces that need text/JSON
// encoding while still being easy to convert back to raw bytes.
package bytes
