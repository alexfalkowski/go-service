// Package cache provides cache abstractions, configuration, and drivers for go-service.
//
// The primary entrypoint is `NewCache`, which constructs a `*Cache` from configuration and registers
// lifecycle hooks to flush/close the underlying cache driver on shutdown.
//
// # Disabled / nil behavior
//
// Caching is intentionally optional. When cache configuration is disabled/unset, constructors return nil
// and callers are expected to tolerate a nil cache instance.
//
// In addition to the instance API on `*Cache`, this package exposes package-level generic helpers
// (`Get` and `Persist`). Those helpers are nil-safe after `Register` has been called (via DI wiring in
// `Module`), and they become no-ops / return zero values when caching is disabled.
//
// # Value encoding
//
// `Cache` persists arbitrary values by encoding (and optionally compressing) them before passing them to
// the configured driver. The encoder/compressor used is selected by configuration with sensible defaults.
package cache
