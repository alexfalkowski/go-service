package json

import (
	"encoding/json"
	"io"
)

// NewEncoder constructs a JSON encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements JSON encoding and decoding.
type Encoder struct{}

// Encode writes v as JSON to w.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode reads JSON from r into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
