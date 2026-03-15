// Package io provides small wrappers and helpers around the standard library io package.
//
// The primary purpose of this package is to offer a stable go-service import path for commonly used
// io primitives while preserving the exact semantics of the standard library `io` package.
//
// Most identifiers in this package are thin aliases/wrappers (for example `Reader`, `Writer`,
// `ReaderFrom`, `WriterTo`, `ReadCloser`, `Resetter`, and `NopCloser`). They exist so packages within
// go-service and downstream services can consistently
// import `github.com/alexfalkowski/go-service/v2/io` without mixing direct stdlib imports across the
// codebase.
//
// This package also provides small convenience helpers, such as `ReadAll`, which captures an entire
// stream and returns both the bytes and a fresh ReadCloser over those bytes for re-reading.
//
// In short, use this package when you want:
//   - stable aliases for common stream interfaces,
//   - a shared `Resetter` abstraction for reusable buffers/readers, or
//   - convenience helpers built on top of standard I/O behavior.
package io
