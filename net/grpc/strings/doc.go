// Package strings provides gRPC-specific string helpers for net-layer middleware.
//
// This package contains helpers for parsing gRPC full method names and matching standard transport
// operation RPCs that selected middleware can treat specially, such as auth, logging, or unary-only limiter
// bypasses.
//
// Start with `SplitServiceMethod` and `IsOperationMethod`.
package strings
