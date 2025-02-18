package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

// NewEncoder for yaml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for yaml.
type Encoder struct{}

// Encode for proto.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

// Decode for proto.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
