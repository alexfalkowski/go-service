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
// The built-in Redis backend resolves its URL from a go-service "source
// string", constructs a go-redis client, and instruments that client via
// `cache/telemetry` before exposing it through the cachego Redis adapter.
// Redis configuration is strict by design: `Config.Options["url"]` must exist
// and be a string. The standard config fixtures provide that shape; callers that
// build config manually should validate it before calling `NewDriver`.
//
// The built-in `sync` driver comes from the upstream cachego dependency and currently has whole-second TTL
// resolution. Callers should not rely on sub-second expiration with that backend.
//
// If the configured kind is unknown, `NewDriver` returns `ErrNotFound`.
//
// # Errors
//
// This package re-exports `cachego.ErrCacheExpired` as `ErrExpired` and provides
// `IsExpiredError` / `IsMissingError` helpers to classify backend-specific miss conditions in a
// backend-agnostic way.
package driver
