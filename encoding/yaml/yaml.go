package yaml

import (
	"io"

	yaml "go.yaml.in/yaml/v3"
)

// NewEncoder constructs a YAML encoder.
//
// This encoder is a thin adapter around go-yaml v3 (imported as go.yaml.in/yaml/v3) that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements YAML encoding and decoding.
//
// It uses go-yaml v3 with default settings.
type Encoder struct{}

// Encode writes v to w as YAML.
//
// This is a thin wrapper around `yaml.NewEncoder(w).Encode(v)`.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

// Decode reads YAML from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// This is a thin wrapper around `yaml.NewDecoder(r).Decode(v)`.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
