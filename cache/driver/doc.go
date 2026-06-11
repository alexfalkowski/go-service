// Package driver provides cache driver construction and related helpers for go-service.
//
// It contains the [NewDriver] constructor used by DI wiring to build a cache backend implementation
// from [github.com/alexfalkowski/go-service/v2/cache/config.Config].
//
// # Disabled / nil behavior
//
// When caching is disabled (i.e. the cache config is nil), [NewDriver] returns a nil [Driver] and a nil error.
//
// # Supported kinds
//
// The driver kind is selected by [github.com/alexfalkowski/go-service/v2/cache/config.Config.Kind].
// Supported values are implementation-dependent, but this package currently includes built-in constructors
// for common backends (for example Redis and an in-memory sync driver).
//
// The built-in Redis backend resolves its URL from a go-service "source
// string", constructs a [github.com/redis/go-redis/v9] client, and instruments that client via
// [github.com/alexfalkowski/go-service/v2/cache/telemetry] before exposing it through a context-aware driver.
// Redis configuration is strict by design: the
// [github.com/alexfalkowski/go-service/v2/cache/config.Config.Options] map must contain a "url" string.
// The standard config fixtures provide that shape; callers that build config manually should validate it
// before calling [NewDriver].
//
// The built-in "sync" driver uses a bounded in-process cache and expires entries
// when they are read or before new values are saved.
//
// If the configured kind is unknown, [NewDriver] returns [ErrNotFound].
//
// # Errors
//
// This package provides [ErrExpired], [ErrMissing], and helper functions to
// classify backend-specific miss conditions in a backend-agnostic way.
package driver
