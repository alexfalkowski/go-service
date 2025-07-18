package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// NewEncoder for toml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for toml.
type Encoder struct{}

// Encode for proto.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

// Decode for proto.
func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)
	return err
}
