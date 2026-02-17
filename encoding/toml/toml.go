package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// NewEncoder constructs a TOML encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements TOML encoding and decoding.
type Encoder struct{}

// Encode writes v as TOML to w.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

// Decode reads TOML from r into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)
	return err
}
