package encoding

import "io"

// Encoder encodes values to a writer and decodes values from a reader.
//
// Encoder is intentionally minimal so multiple concrete encodings (JSON/YAML/TOML/protobuf/gob, etc.)
// can be used interchangeably.
//
// # Encode contract
//
// Encode must serialize v to w. Implementations may require that v satisfies additional interfaces
// or is of a particular shape (for example a protobuf encoder may require v to implement
// google.golang.org/protobuf/proto.Message).
//
// # Decode contract
//
// Decode must read from r and populate v. In most cases v is expected to be a pointer to the target
// value so the decoder can mutate it (e.g. *MyStruct). Implementations may return an error if v is not
// a supported type (for example encoding/errors.ErrInvalidType).
//
// Implementations should return any underlying I/O errors and any parse/unmarshal errors produced by
// their respective codecs.
type Encoder interface {
	// Encode writes a serialized representation of v to w.
	Encode(w io.Writer, v any) error

	// Decode reads from r and decodes into v.
	Decode(r io.Reader, v any) error
}
