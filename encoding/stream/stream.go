package stream

import (
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
)

// Encoder encodes a sequence of values to a writer bound at construction.
//
// Encoder differs from [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder] in two ways: it is
// bound to one writer for its entire lifetime (constructed with the writer rather than given one per
// call), and it is expected to serve many Encode calls rather than exactly one.
//
// # Encode contract
//
// Encode serializes v and writes it to the underlying writer. Implementations may require that v
// satisfies additional interfaces or is of a particular shape, matching the wrapped codec's own
// requirements. Callers may call Encode repeatedly to write multiple values to the same stream.
//
// # Close contract
//
// Close finalizes any codec-level state (for example a stream terminator or trailer). Close must
// never close the underlying writer — the writer is owned by the caller, not the Encoder. Close is
// called exactly once, at the true end of the stream: implementations are not required to tolerate a
// second call (for example encoding/stream/yaml's Encoder is a direct alias of a non-idempotent
// upstream Close), so callers must not both defer Close and call it explicitly on the success path.
type Encoder interface {
	// Encode writes a serialized representation of v to the underlying writer.
	Encode(v any) error

	// Close finalizes the encoder. Close never closes the underlying writer. Close is called exactly
	// once; a second call is not guaranteed to be safe.
	Close() error
}

// Decoder decodes a sequence of values from a reader bound at construction.
//
// Decoder differs from [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder]'s Decode in two
// ways: it is bound to one reader for its entire lifetime (constructed with the reader rather than
// given one per call), and it is expected to serve many Decode calls rather than exactly one.
//
// # Decode contract
//
// Decode reads the next encoded value from the underlying reader and populates v. In most cases v
// should be a pointer to the destination value (for example *MyStruct). Decode returns io.EOF once
// the underlying stream is exhausted, matching the terminal behavior of the wrapped codec's own
// decoder.
//
// # Close contract
//
// Close finalizes any codec-level state. Close must never close the underlying reader — the reader
// is owned by the caller, not the Decoder. Close is called exactly once; implementations are not
// required to tolerate a second call.
type Decoder interface {
	// Decode reads the next value from the underlying reader into v. Decode returns io.EOF once the
	// stream is exhausted.
	Decode(v any) error

	// Close finalizes the decoder. Close never closes the underlying reader. Close is called exactly
	// once; a second call is not guaranteed to be safe.
	Close() error
}

// EncoderFunc constructs an [Encoder] bound to w. Registered in a [Map] by kind as part of a [Codec]
// via [Map.Register], and returned by [Map.Get].
type EncoderFunc func(w io.Writer) Encoder

// DecoderFunc constructs a [Decoder] bound to r. Registered in a [Map] by kind as part of a [Codec]
// via [Map.Register], and returned by [Map.Get].
//
// With [codec.WithDiscardUnknown], the decoder ignores source members that
// do not map to its destination.
type DecoderFunc func(r io.Reader, opts ...codec.Option) Decoder
