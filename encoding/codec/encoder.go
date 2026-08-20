// Package codec defines the common contract implemented by supported encodings.
package codec

import "github.com/alexfalkowski/go-service/v2/io"

// Encoder encodes values to a writer and decodes values from a reader.
//
// Encoder is intentionally minimal so multiple concrete encodings (JSON/YAML/TOML/protobuf/gob, etc.)
// can be used interchangeably.
//
// # Encode contract
//
// Encode must serialize v to w. Implementations may require that v satisfies additional interfaces
// or is of a particular shape (for example a protobuf encoder may require v to implement
// google.golang.org/protobuf/proto.Message). Implementations should reject a value with missing
// required fields when their format supports that distinction, unless WithAllowPartial is supplied.
//
// # Decode contract
//
// Decode must read from r and populate v. In most cases v is expected to be a pointer to the target
// value so the decoder can mutate it (e.g. *MyStruct). Implementations may return an error if v is not
// a supported type (for example [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType]).
//
// Structured single-value decoders should reject additional encoded values after the first payload,
// either by consuming the whole input or by returning [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrTrailingData]. Stream or
// passthrough encoders may delegate full-consumption semantics to the concrete value they decode into.
//
// Some implementations buffer the remaining contents of r before decoding. When r contains untrusted
// input, callers must bound it before calling Decode. Standard go-service HTTP and cache wiring applies
// those limits before values reach encoders.
//
// Implementations should return any underlying I/O errors and any parse/unmarshal errors produced by
// their respective codecs.
type Encoder interface {
	// Encode writes a serialized representation of v to w. WithAllowPartial lets codecs that
	// validate required fields (for example protobuf) encode a message with missing required
	// fields instead of returning an error; codecs without that concept ignore it.
	Encode(w io.Writer, v any, opts ...Option) error

	// Decode reads from r and decodes into v. WithDiscardUnknown makes codecs
	// ignore source members that do not map to v; without it, codecs reject
	// unknown members when their format supports that distinction. WithAllowPartial
	// makes codecs that validate required fields (for example protobuf) decode a
	// value with missing required fields instead of returning an error; codecs
	// without that concept ignore it. Both may be supplied together.
	Decode(r io.Reader, v any, opts ...Option) error
}
