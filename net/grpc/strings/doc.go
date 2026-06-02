// Package strings provides gRPC-specific string helpers for net-layer middleware.
//
// This package contains helpers for matching standard transport operation RPCs that selected middleware can
// treat specially, such as auth, logging, or unary-only limiter bypasses.
//
// Start with [IsOperationMethod].
package strings
