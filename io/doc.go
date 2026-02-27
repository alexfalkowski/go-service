// Package io provides small wrappers and helpers around the standard library io package.
//
// The primary purpose of this package is to offer a stable go-service import path for commonly used
// io primitives while preserving the exact semantics of the standard library `io` package.
//
// Most identifiers in this package are thin aliases/wrappers (for example `Reader`, `ReadCloser`,
// and `NopCloser`). They exist so packages within go-service and downstream services can consistently
// import `github.com/alexfalkowski/go-service/v2/io` without mixing direct stdlib imports across the
// codebase.
//
// This package also provides small convenience helpers, such as `ReadAll`, which captures an entire
// stream and returns both the bytes and a fresh ReadCloser over those bytes for re-reading.
package io
