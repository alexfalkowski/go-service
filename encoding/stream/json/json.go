package json

import (
	"encoding/json"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
)

// Encoder streams JSON values to a writer bound at construction.
//
// Encoder embeds the standard library's [encoding/json.Encoder], so [Encoder.Encode] is the embedded
// encoder's own method. Encoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder] structurally, without this
// package importing that package: [github.com/alexfalkowski/go-service/v2/encoding/stream.NewMap]
// imports this package to pre-register the default "json" kind, and this package importing that one
// back would be a compile-time import cycle.
type Encoder struct {
	*json.Encoder
}

// NewEncoder constructs a streaming JSON [Encoder] bound to w.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/json.Encoder], the returned encoder writes
// no indentation, so consecutive values stay newline-delimited (NDJSON) instead of each spanning
// multiple lines.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Encoder: json.NewEncoder(w)}
}

// Close is a no-op: the standard library JSON encoder has no finalization state to flush.
func (e *Encoder) Close() error {
	return nil
}

// Decoder streams JSON values from a reader bound at construction.
//
// Decoder embeds the standard library's [encoding/json.Decoder], so [Decoder.Decode] is the embedded
// decoder's own method. Decoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] structurally, without this
// package importing that package.
type Decoder struct {
	*json.Decoder
}

// NewDecoder constructs a streaming JSON [Decoder] bound to r.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/json.Encoder]'s Decode, the returned
// decoder allows additional JSON values after the first: it does not enforce single-value
// trailing-data rejection, since a stream is expected to carry many values. Unknown fields are
// rejected unless [codec.WithDiscardUnknown] is supplied.
func NewDecoder(r io.Reader, opts ...codec.Option) *Decoder {
	resolved := codec.Apply(opts...)

	decoder := json.NewDecoder(r)
	if !resolved.DiscardUnknown() {
		decoder.DisallowUnknownFields()
	}

	return &Decoder{Decoder: decoder}
}

// Close is a no-op: the standard library JSON decoder has no finalization state to flush.
func (d *Decoder) Close() error {
	return nil
}
