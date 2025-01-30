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

func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)

	return err
}
