package json

import (
	"encoding/json"
	"io"
)

// NewEncoder constructs a JSON encoder.
//
// This encoder is a thin adapter around the standard library `encoding/json` package that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements JSON encoding and decoding.
//
// It uses the standard library `encoding/json` encoder/decoder with default settings.
type Encoder struct{}

// Encode writes v to w as JSON.
//
// This is a thin wrapper around `json.NewEncoder(w).Encode(v)`.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode reads JSON from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// This is a thin wrapper around `json.NewDecoder(r).Decode(v)`.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
