// Package strings provides gRPC-specific string helpers for net-layer middleware.
//
// This package contains helpers for parsing gRPC full method names and matching standard operational
// RPCs that are commonly ignored by auth, limiter, and logging middleware.
//
// Start with `SplitServiceMethod` and `IsIgnorable`.
package strings
