// Package encoding provides value encoding/decoding helpers and DI wiring used by go-service.
//
// This package provides a registry ([Map]) that selects
// [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder] implementations by kind at runtime.
//
// # Registry
//
// Map is a kind-to-[github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder] lookup. It is commonly used by configuration loading and transport layers to
// choose a decoder/encoder based on either:
//   - a file extension (for example "yaml", "toml", "json"), or
//   - a content kind / media subtype (for example "protobuf", "bytes").
//
// Map registers each encoder under exactly one canonical kind; callers that need to accept alternate
// spellings (such as HTTP media subtype aliases "pb" or "octet-stream") translate them to the canonical
// kind before calling [Map.Get] rather than relying on this registry to know every alias.
//
// Callers typically obtain a *[Map] via DI and then use [Map.Get] to select an encoder, often
// falling back to a default when the requested kind is not registered.
//
// # Decode options
//
// Decoding rejects unknown members by default where the format supports that distinction. Callers at
// a forward-compatible API boundary can pass
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] to
// [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder.Decode]; framing,
// syntax, and type validation remain unchanged.
//
// # Wiring
//
// NewMap constructs a *[Map] that registers default encoders under common kinds used throughout
// go-service, including:
//   - JSON, HJSON, YAML, TOML, MessagePack
//   - protobuf binary/text/JSON variants
//   - gob
//   - "bytes" passthrough for [io.ReaderFrom]/[io.WriterTo] payloads
//
// Module provides the default *[Map] for Fx applications.
//
// Start with [Map], [NewMap], [Module], and
// [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder].
package encoding
