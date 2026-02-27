// Package sync provides synchronization primitives and pooling helpers used by go-service.
//
// This package intentionally offers a small, stable API surface that wraps:
//
//   - The standard library sync primitives (for example Mutex and RWMutex).
//   - A small set of additional concurrency-safe utilities from the external
//     module github.com/alexfalkowski/go-sync (for example a generic Map and a
//     byte buffer pool).
//
// The goal is to let go-service code depend on go-service packages consistently,
// while still using well-known synchronization patterns and implementations.
//
// # Standard library wrappers
//
// This package re-exports common mutex types as aliases:
//
//   - Mutex (alias of sync.Mutex)
//   - RWMutex (alias of sync.RWMutex)
//
// Because these are type aliases, they have the same semantics and methods as
// the standard library types.
//
// # Generic concurrent map
//
// Map and NewMap are aliases/wrappers around github.com/alexfalkowski/go-sync's
// generic concurrent map implementation:
//
//   - Map[K, V] is a concurrent map keyed by K with values V.
//   - NewMap[K, V] constructs a new instance.
//
// Consult the upstream go-sync documentation for details of the map’s method set
// and concurrency guarantees.
//
// # Buffer pool
//
// BufferPool and NewBufferPool are aliases/wrappers around github.com/alexfalkowski/go-sync's
// buffer pool implementation:
//
//   - BufferPool provides pooled buffers to reduce allocations in hot paths.
//   - NewBufferPool constructs a new pool.
//
// Buffer pools are typically used by transports/encoders to reuse temporary
// buffers when assembling payloads.
//
// # Dependency injection (Fx)
//
// This package also exports Module, which wires the buffer pool into Fx.
// Including Module in an Fx application provides a shared *BufferPool instance
// constructed via the upstream go-sync constructor.
//
// # Notes
//
// This package does not attempt to replace the standard library sync package.
// If you need primitives not exposed here, prefer importing sync directly.
package sync
