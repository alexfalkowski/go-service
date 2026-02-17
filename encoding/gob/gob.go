package gob

import (
	"encoding/gob"
	"io"
)

// NewEncoder constructs a gob encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements gob encoding and decoding.
type Encoder struct{}

// Encode writes v as gob to w.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

// Decode reads gob from r into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}
