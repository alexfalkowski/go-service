package msgpack

import (
	"github.com/Basekick-Labs/msgpack/v6"
	"github.com/alexfalkowski/go-service/v2/io"
)

// Encoder streams MessagePack values to a writer bound at construction.
//
// Encoder embeds the [github.com/Basekick-Labs/msgpack/v6.Encoder], so [Encoder.Encode] is the
// embedded encoder's own method. Encoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Encoder] structurally, without this
// package importing that package: [github.com/alexfalkowski/go-service/v2/encoding/stream.NewMap]
// imports this package to pre-register the default "msgpack" kind, and this package importing that
// one back would be a compile-time import cycle.
type Encoder struct {
	*msgpack.Encoder
}

// NewEncoder constructs a streaming MessagePack [Encoder] bound to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Encoder: msgpack.NewEncoder(w)}
}

// Close is a no-op: the underlying MessagePack encoder has no finalization state to flush.
func (e *Encoder) Close() error {
	return nil
}

// Decoder streams MessagePack values from a reader bound at construction.
//
// Decoder embeds the [github.com/Basekick-Labs/msgpack/v6.Decoder], so [Decoder.Decode] is the
// embedded decoder's own method. Decoder satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/stream.Decoder] structurally, without this
// package importing that package.
type Decoder struct {
	*msgpack.Decoder
}

// NewDecoder constructs a streaming MessagePack [Decoder] bound to r.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding/msgpack.Encoder]'s Decode, the returned
// decoder allows additional MessagePack values after the first: it does not enforce single-value
// trailing-data rejection, since a stream is expected to carry many values.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{Decoder: msgpack.NewDecoder(r)}
}

// Close is a no-op: the underlying MessagePack decoder has no finalization state to flush.
func (d *Decoder) Close() error {
	return nil
}
