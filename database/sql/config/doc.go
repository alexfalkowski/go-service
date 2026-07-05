// Package config provides shared SQL ([github.com/alexfalkowski/go-service/v2/database/sql]) configuration types for go-service.
//
// This package defines configuration structs that are typically embedded into a larger service
// configuration and then consumed by driver-specific wiring (for example [github.com/alexfalkowski/go-service/v2/database/sql/pg]).
//
// # Pool configuration
//
// [Config] models reader and writer [Pool] values. Each pool owns both its
// datasource names and its complete [PoolSettings], including:
//   - maximum connection lifetime,
//   - max open connections, and
//   - max idle connections.
//
// # DSNs and source strings
//
// Reader and writer DSNs are configured through [Pool.DSNs]. Each DSN URL is expressed as a
// go-service "source string" (resolved via [github.com/alexfalkowski/go-service/v2/os.FS.ReadSource]), so it can be:
//   - "env:NAME" to read from an environment variable,
//   - "file:/path/to/dsn" to read from a file, or
//   - any other value treated as a literal DSN.
//
// Enabled SQL configurations must provide at least one writer or reader DSN, and each resolved DSN must be non-empty.
//
// Start with [Config], [Pool], [PoolSettings], and [DSN].
package config
