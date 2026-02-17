package encoding

import "io"

// Encoder encodes values to a writer and decodes values from a reader.
type Encoder interface {
	// Encode writes v to w.
	Encode(w io.Writer, v any) error

	// Decode reads from r into v.
	Decode(r io.Reader, v any) error
}
