package gob

import (
	"encoding/gob"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

// NewEncoder constructs a gob encoder.
//
// This encoder is a thin adapter around the standard library `encoding/gob` package that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements gob encoding and decoding.
//
// It uses the standard library `encoding/gob` encoder/decoder with default settings.
type Encoder struct{}

// Encode writes v to w as gob.
//
// This is a thin wrapper around `gob.NewEncoder(w).Encode(v)`.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

// Decode reads gob from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// Decode reads one gob value and rejects additional values in the same stream.
func (e *Encoder) Decode(r io.Reader, v any) error {
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	var extra any
	return errors.TrailingData(decoder.Decode(&extra))
}
