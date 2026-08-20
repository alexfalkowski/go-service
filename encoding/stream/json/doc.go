// Package json provides a streaming JSON [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]/
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] pair used by go-service.
//
// It wraps the standard library [encoding/json] package's own NewEncoder/NewDecoder types, which are
// each bound to one writer/reader for many Encode/Decode calls. It rejects unknown fields by default
// and accepts [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] for
// forward-compatible API streams. It drops the sibling package's indentation (which would break
// newline-delimited framing) and trailing-data rejection (a stream is expected to contain many values).
//
// Start with the package-level constructors.
package json
