package yaml

import (
	"io"

	yaml "go.yaml.in/yaml/v2"
)

// NewEncoder for yaml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for yaml.
type Encoder struct{}

// Encode for yaml.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

// Decode for yaml.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
