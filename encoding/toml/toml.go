package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// Encoder for toml.
type Encoder struct{}

// NewEncoder for toml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)

	return err
}
