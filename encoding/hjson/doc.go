// Package hjson provides HJSON encoding helpers and adapters used by go-service.
//
// This package integrates HJSON encoding/decoding behind the go-service encoding abstraction.
// Decode rejects duplicate object keys and unknown destination fields by default; callers can use
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] for forward-compatible API
// payloads.
//
// Start with the package-level constructors.
package hjson
