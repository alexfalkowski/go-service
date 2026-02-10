// Package hmac provides HMAC signing and verification helpers for go-service.
//
// This package wires an HMAC-based Signer that uses HMAC-SHA-512. Keys are typically loaded via the
// go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value) using os.FS.ReadSource.
//
// Verification uses crypto/hmac.Equal for constant-time MAC comparison.
package hmac
