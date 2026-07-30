// Package json provides a streaming JSON [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]/
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] pair used by go-service.
//
// It wraps the standard library [encoding/json] package's own NewEncoder/NewDecoder types, which are
// each bound to one writer/reader for many Encode/Decode calls. It carries over the strict decoding
// policy of [github.com/alexfalkowski/go-service/v2/encoding/json] (unknown fields rejected), but
// drops that package's indentation (which would break newline-delimited framing) and trailing-data
// rejection (a stream is expected to contain many values).
//
// Start with the package-level constructors.
package json
