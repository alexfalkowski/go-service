package json

import (
	"encoding/json"
	"io"
)

// NewEncoder for json.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for json.
type Encoder struct{}

// Encode for json.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode for json.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
