// Package io provides small wrappers and helpers around the standard library io package.
//
// The primary purpose of this package is to offer a stable go-service import path for commonly used
// io primitives.
//
// Most identifiers in this package are thin standard-library aliases/wrappers (for example `Reader`,
// `Writer`, `ReaderFrom`, `WriterTo`, `ReadCloser`, `NopCloser`, and `WriteString`). They preserve
// standard-library semantics while letting packages within go-service and downstream services consistently
// import `github.com/alexfalkowski/go-service/v2/io` without mixing direct stdlib imports across the
// codebase.
//
// This package also provides small repository-specific helpers. `Resetter` is a local abstraction
// for values that can be reset before reuse, and `ReadAll` captures an entire stream and returns
// both the bytes and a fresh ReadCloser over those bytes for re-reading.
//
// In short, use this package when you want:
//   - stable aliases for common stream interfaces,
//   - a shared `Resetter` abstraction for reusable buffers/readers, or
//   - convenience helpers built on top of standard I/O behavior.
package io
