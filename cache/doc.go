// Package cache provides cache abstractions, configuration, and drivers for go-service.
//
// The primary entrypoint is [NewCache], which constructs a *[Cache] from configuration.
//
// # Disabled / nil behavior
//
// Caching is intentionally optional. When cache configuration is disabled/unset, constructors return nil
// and callers are expected to tolerate a nil cache instance.
//
// In addition to the instance API on *[Cache], this package exposes package-level generic helpers
// ([Get] and [Persist]). Those helpers are nil-safe after [Register] has been called (via DI wiring in
// [Module]), and they become no-ops / return zero values when caching is disabled. In the standard
// service composition this registration is performed for you by the module graph.
//
// # Value encoding
//
// [Cache] persists arbitrary values by encoding (and optionally compressing) them before passing them to
// the configured driver. The encoder/compressor used is selected by configuration with sensible defaults.
// The configured encoder and compressor are also included in the driver key namespace so format changes
// create cache misses instead of decoding values written by an incompatible format.
//
// # TTL resolution
//
// TTL handling depends on the selected driver. The built-in in-memory "ttlcache"
// driver stores a bounded number of values in process memory, expires entries
// when they are read, and removes expired entries before saving new values.
//
// # Flush behavior
//
// [Cache.Flush] delegates to the selected driver and can have backend-wide
// effects. The built-in Redis backend uses Redis FLUSHDB, so it clears the
// entire selected Redis database, including keys that were not created through
// this cache facade. Use a dedicated Redis database for go-service cache data
// before calling Flush against Redis.
package cache
