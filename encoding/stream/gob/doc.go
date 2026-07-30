// Package gob provides a streaming gob [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]/
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] pair used by go-service.
//
// It wraps the standard library [encoding/gob] package's own NewEncoder/NewDecoder types with one
// persistent encoder/decoder bound to the writer/reader for the lifetime of the stream, unlike
// [github.com/alexfalkowski/go-service/v2/encoding/gob], which constructs and discards a fresh
// encoder/decoder per call and rejects trailing values because it models exactly one value per call.
//
// Start with the package-level constructors.
package gob
