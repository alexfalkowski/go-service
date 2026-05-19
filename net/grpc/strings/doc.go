// Package strings provides gRPC-specific string helpers for net-layer middleware.
//
// This package contains helpers for parsing gRPC full method names and matching standard transport
// operation RPCs that can bypass auth, limiter, and logging middleware.
//
// Start with `SplitServiceMethod` and `IsOperationMethod`.
package strings
