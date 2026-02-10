// Package pem provides PEM decoding helpers used by go-service.
//
// This package provides a Decoder that reads PEM-encoded data (typically via the go-service "source string"
// pattern such as "env:NAME", "file:/path", or a literal value) and returns the raw bytes of a specific PEM block kind.
//
// Start with `Decoder` and `NewDecoder`.
package pem
