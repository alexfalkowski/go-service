package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

// Encoder for yaml.
type Encoder struct{}

// NewEncoder for yaml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
