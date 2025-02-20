package encoding

import "io"

// Encoder allows different types of encoding/decoding.
type Encoder interface {
	// Encode any to a writer.
	Encode(w io.Writer, e any) error

	// Decode any from a reader.
	Decode(r io.Reader, e any) error
}
