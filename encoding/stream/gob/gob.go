package gob

import (
	"encoding/gob"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
)

// Encoder streams gob values to a writer bound at construction.
//
// Encoder embeds the standard library's [encoding/gob.Encoder], so [Encoder.Encode] is the embedded
// encoder's own method. Encoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder] structurally, without this
// package importing that package: [github.com/alexfalkowski/go-service/v2/encoding/stream.NewMap]
// imports this package to pre-register the default "gob" kind, and this package importing that one
// back would be a compile-time import cycle.
type Encoder struct {
	*gob.Encoder
}

// NewEncoder constructs a streaming gob [Encoder] bound to w.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/gob.Encoder], which constructs a fresh
// standard library encoder per Encode call, the returned encoder wraps one persistent
// [encoding/gob.Encoder] bound to w for the lifetime of the stream — the intended way to use gob for
// more than one value on the same connection.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Encoder: gob.NewEncoder(w)}
}

// Close is a no-op: the standard library gob encoder has no finalization state to flush.
func (e *Encoder) Close() error {
	return nil
}

// Decoder streams gob values from a reader bound at construction.
//
// Decoder embeds the standard library's [encoding/gob.Decoder], so [Decoder.Decode] is the embedded
// decoder's own method. Decoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] structurally, without this
// package importing that package.
type Decoder struct {
	*gob.Decoder
}

// NewDecoder constructs a streaming gob [Decoder] bound to r.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/gob.Encoder]'s Decode, the returned decoder
// allows additional gob values after the first: it does not enforce single-value trailing-data
// rejection, since a stream is expected to carry many values.
func NewDecoder(r io.Reader, _ ...codec.Option) *Decoder {
	return &Decoder{Decoder: gob.NewDecoder(r)}
}

// Close is a no-op: the standard library gob decoder has no finalization state to flush.
func (d *Decoder) Close() error {
	return nil
}
