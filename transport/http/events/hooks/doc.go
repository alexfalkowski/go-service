// Package hooks provides CloudEvents-specific HTTP webhook middleware for go-service.
//
// This package re-exports the generic HTTP webhook transport wrapper under a CloudEvents-focused import path
// and provides a small adapter that wraps a standard `http.Handler`.
//
// Disabled behavior:
// When webhook support is disabled, the returned handler behaves as a pass-through and simply delegates to
// the wrapped handler.
package hooks
