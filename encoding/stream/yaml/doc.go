// Package yaml provides a streaming YAML [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]/
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] pair used by go-service.
//
// It wraps [go.yaml.in/yaml/v3]'s own NewEncoder/NewDecoder types, which are each bound to one
// writer/reader for many Encode/Decode calls. It carries over the strict decoding policy of
// [github.com/alexfalkowski/go-service/v2/encoding/yaml] (unknown fields rejected), but drops that
// package's trailing-document rejection, since a stream is expected to contain many documents. Encoder
// is a direct alias of the upstream type, so its Close is not idempotent — a second call surfaces
// upstream's "expected nothing after STREAM-END" error; callers must call Close exactly once, at the
// true end of the stream, per [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]'s
// contract. See [NewEncoder].
//
// Start with the package-level constructors.
package yaml
