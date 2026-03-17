// Package strings provides HTTP-specific string helpers for net-layer middleware.
//
// This package contains small helpers used by HTTP middleware and telemetry code, including matching
// well-known operational endpoints such as health and metrics routes that are typically ignored by
// auth, limiter, and logging middleware.
//
// Start with `IsIgnorable`.
package strings
