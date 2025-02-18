package gob

import (
	"encoding/gob"
	"io"
)

// NewEncoder for gob.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for gob.
type Encoder struct{}

// Encode for gob.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

// Decode for gob.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}
