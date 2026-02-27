package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// NewEncoder constructs a TOML encoder.
//
// This encoder is a thin adapter around github.com/BurntSushi/toml that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements TOML encoding and decoding.
//
// It uses BurntSushi/toml with default settings.
type Encoder struct{}

// Encode writes v to w as TOML.
//
// This is a thin wrapper around `toml.NewEncoder(w).Encode(v)`.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

// Decode reads TOML from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
//
// This method intentionally discards metadata returned by BurntSushi/toml and returns only the
// decode/unmarshal error (if any).
func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)
	return err
}
