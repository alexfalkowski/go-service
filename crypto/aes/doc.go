// Package aes provides AES-GCM encryption helpers and wiring for go-service.
//
// This package wires an AES-GCM Cipher that encrypts and decrypts byte slices using a configured key.
// Keys are typically loaded via the go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value).
package aes
