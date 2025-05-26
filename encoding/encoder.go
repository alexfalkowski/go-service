package encoding

import "io"

// Encoder allows different types of encoding/decoding.
type Encoder interface {
	// Encode any to a writer.
	Encode(w io.Writer, v any) error

	// Decode any from a reader.
	Decode(r io.Reader, v any) error
}
