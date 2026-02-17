package yaml

import (
	"io"

	yaml "go.yaml.in/yaml/v3"
)

// NewEncoder constructs a YAML encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements YAML encoding and decoding.
type Encoder struct{}

// Encode writes v as YAML to w.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

// Decode reads YAML from r into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
