// Package errors provides small error helpers used across go-service.
//
// This package is intentionally lightweight. It primarily re-exports a subset of the standard library
// `errors` package APIs (As/Is/Join/New) behind a stable go-service import path, and provides a small
// convenience helper (`Prefix`) for consistently attributing errors to a subsystem or component.
//
// # Re-exports
//
// The following functions mirror the behavior and semantics of the standard library equivalents:
//
//   - As: type-assert (via error chain traversal) into a target
//   - Is: match a target error in an error chain
//   - Join: combine multiple errors into one
//   - New: construct a sentinel error value
//
// # Prefixing errors
//
// Prefix is commonly used at module/service boundaries to add a stable component prefix to an error
// message while preserving the original error for unwrapping:
//
//	err := errors.Prefix("cache", underlyingErr)
//
// Prefix returns nil when the input error is nil, which makes it convenient to use in return
// statements without additional nil checks.
//
// Start with `As`, `Is`, `Join`, `New`, and `Prefix`.
package errors
