// Package net provides small network helpers and wrappers used by go-service.
//
// This package is a lightweight convenience layer over the standard library `net` package.
// It exists to provide a stable go-service import path for common networking primitives and
// a few helpers that are used by higher-level server and transport wiring.
//
// # Re-exported types
//
// This package re-exports selected `net` types (Conn, Dialer, Listener) as aliases so callers can
// depend on go-service primitives while preserving the exact semantics of the standard library.
//
// # Helpers
//
// The helpers in this package focus on:
//   - creating listeners that respect context cancellation (Listen), and
//   - working with go-service address conventions, such as "tcp://host:port" style addresses
//     (SplitNetworkAddress, Host, DefaultAddress).
//
// Start with `Listen` and `DefaultAddress`.
package net
