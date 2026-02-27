// Package encoding provides value encoding/decoding helpers and DI wiring used by go-service.
//
// This package defines a small `Encoder` interface (encode to an `io.Writer`, decode from an `io.Reader`)
// and a registry (`Map`) used to select an encoder by kind at runtime.
//
// # Registry
//
// `Map` is a kind→Encoder lookup. It is commonly used by configuration loading and transport layers to
// choose a decoder/encoder based on either:
//   - a file extension (e.g. "yaml", "yml", "toml", "json"), or
//   - a content kind / media subtype (e.g. "proto", "plain", "octet-stream").
//
// Callers typically obtain a `*Map` via DI and then use `Get(kind)` to select an encoder, often
// falling back to a default when the requested kind is not registered.
//
// # Wiring
//
// `Module` wires the default encoder implementations and provides a `*Map` pre-populated with common
// kinds used throughout go-service, including:
//   - JSON, YAML, TOML
//   - protobuf binary/text/JSON variants
//   - gob
//   - "plain"/bytes passthrough for io.ReaderFrom/io.WriterTo payloads
//
// Start with `Encoder`, `Map`, `NewMap`, and `Module`.
package encoding
