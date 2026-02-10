// Package crypto provides cryptographic helpers and wiring used by go-service.
//
// This package groups cryptographic primitives and integrations that are used across
// the module, such as hashing/signing, encryption, key parsing, and TLS helpers.
//
// Prefer the concrete subpackages (for example `crypto/tls`, `crypto/rsa`, `crypto/hmac`)
// for specific functionality.
package crypto
