// Package driver provides cache driver construction and related helpers for go-service.
//
// It contains the `NewDriver` constructor used by DI wiring to build a cache backend implementation
// from `cache/config.Config`.
//
// # Disabled / nil behavior
//
// When caching is disabled (i.e. the cache config is nil), `NewDriver` returns a nil Driver and a nil error.
//
// # Supported kinds
//
// The driver kind is selected by `Config.Kind`. Supported values are implementation-dependent, but this
// package currently includes built-in constructors for common backends (for example Redis and an in-memory
// sync driver).
//
// If the configured kind is unknown, `NewDriver` returns `ErrNotFound`.
//
// # Errors
//
// This package re-exports `cachego.ErrCacheExpired` as `ErrExpired` and provides `IsExpiredError` to
// classify expired-entry errors in a backend-agnostic way.
package driver
