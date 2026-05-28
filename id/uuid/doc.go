// Package uuid provides UUID-based ID generation helpers used by go-service.
//
// This package integrates UUID generation behind the go-service ID abstraction.
//
// # Random pool
//
// Importing this package enables google/uuid's process-wide random pool. UUIDv7 generation sits on
// request metadata hot paths, and the pool is an intentional performance tradeoff for operational
// identifiers. These UUID values are not secrets or bearer tokens.
//
// Start with the package-level constructors.
package uuid
