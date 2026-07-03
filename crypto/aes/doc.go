// Package aes provides AES-GCM encryption helpers and wiring for go-service.
//
// This package wires an AES-GCM Cipher that encrypts and decrypts message data using a configured key.
// Message metadata is authenticated as AES-GCM associated data and must match during decryption.
// Keys are typically loaded via the go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value).
package aes
