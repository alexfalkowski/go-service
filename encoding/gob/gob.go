package gob

import (
	"encoding/gob"
	"io"
)

// Encoder for gob.
type Encoder struct{}

// NewEncoder for gob.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}
