// Package id provides ID generation abstractions, registries, and wiring used by go-service.
//
// This package defines a small `Generator` interface and provides a registry (`Map`) of generators
// keyed by kind (e.g. "uuid", "ksuid", etc.). A concrete generator can be selected at runtime via
// `NewGenerator` using `Config.Kind`.
//
// # Kinds and implementations
//
// Concrete generator implementations live in subpackages under `id/*` (for example `id/uuid`,
// `id/ksuid`, `id/ulid`, `id/nanoid`, and `id/xid`). The `id.Module` wiring constructs these
// generators and registers them into a `*Map`.
//
// # Configuration and enablement
//
// ID generation configuration is optional. By convention across go-service config types, a nil
// `*id.Config` is treated as "disabled" and `NewGenerator` returns (nil, nil) when disabled.
// If configuration is enabled but the configured kind is unknown, `NewGenerator` returns ErrNotFound.
//
// Start with `Generator`, `Config`, `NewGenerator`, `Map`, and `Module`.
package id
