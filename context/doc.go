// Package context provides small wrappers and aliases around the standard library context package.
//
// The primary purpose of this package is to offer a stable go-service import path for commonly used
// context primitives while preserving the exact semantics of the standard library `context` package.
//
// Most identifiers are thin aliases/wrappers around `context` (for example `Context`, `CancelFunc`,
// `CancelCauseFunc`, `Background`, `WithDeadline`, `WithTimeout`, `Cause`, and `WithValue`). They exist so
// packages within go-service and downstream services can consistently import
// `github.com/alexfalkowski/go-service/v2/context` without mixing direct stdlib imports across the codebase.
//
// Cause-aware APIs mirror the standard library behavior: `Context.Err()` continues to report the standard
// sentinels (`Canceled` or `DeadlineExceeded`), while `Cause(ctx)` can expose a richer diagnostic error when
// a cause-aware constructor was used.
//
// This package also defines `Key`, a typed helper for context value keys to reduce accidental collisions
// when multiple packages store values in the same context.
package context
