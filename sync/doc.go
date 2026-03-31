// Package sync wires the shared buffer pool used by go-service.
//
// This package now has a narrow scope: it does not re-export mutexes, maps, or
// buffer pool types. Code that needs synchronization primitives or
// concurrency-safe containers should import the standard library sync package or
// github.com/alexfalkowski/go-sync directly.
//
// # Dependency injection (Fx)
//
// Module registers github.com/alexfalkowski/go-sync.NewBufferPool with the Fx
// graph so packages can share a single upstream buffer pool instance.
//
// This shared pool is consumed by allocation-sensitive helpers such as cache
// encoding and HTTP request/response helpers.
//
// # When to import this package
//
// Import this package when composing go-service Fx modules and you want the
// standard buffer pool wiring.
//
// If you need the concrete pool type or other synchronization helpers, import
// github.com/alexfalkowski/go-sync directly.
package sync
