package yaml

import (
	"github.com/alexfalkowski/go-service/v2/io"
	yaml "go.yaml.in/yaml/v3"
)

// Encoder is an alias of [go.yaml.in/yaml/v3.Encoder].
//
// Encoder already satisfies [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder]
// structurally (Encode and Close), so no wrapper type is needed.
//
// Note: the underlying Close is not idempotent — a second call returns
// "yaml: expected nothing after STREAM-END". Callers must call Close exactly once, at the end of the
// stream.
type Encoder = yaml.Encoder

// NewEncoder constructs a streaming YAML [Encoder] bound to w.
//
// The underlying encoder separates successive documents with "---" and only finalizes the stream once
// Close is called.
func NewEncoder(w io.Writer) *Encoder {
	return yaml.NewEncoder(w)
}

// Decoder streams YAML values from a reader bound at construction.
//
// Decoder embeds the [go.yaml.in/yaml/v3.Decoder], so [Decoder.Decode] is the embedded decoder's own
// method. Decoder satisfies [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder]
// structurally, without this package importing that package: the underlying decoder has no Close
// method at all, so this thin wrapper exists only to add one.
type Decoder struct {
	*yaml.Decoder
}

// NewDecoder constructs a streaming YAML [Decoder] bound to r.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/yaml.Encoder]'s Decode, the returned
// decoder allows additional YAML documents after the first: it does not enforce single-value
// trailing-data rejection, since a stream is expected to carry many documents. Unknown fields remain
// rejected.
func NewDecoder(r io.Reader) *Decoder {
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	return &Decoder{Decoder: decoder}
}

// Close is a no-op: [go.yaml.in/yaml/v3]'s Decoder has no Close method or finalization state.
func (d *Decoder) Close() error {
	return nil
}
