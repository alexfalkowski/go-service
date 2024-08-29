package json

import (
	"io"

	"github.com/goccy/go-json"
)

// Encoder for json.
type Encoder struct{}

// NewEncoder for json.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
