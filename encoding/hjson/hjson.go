package hjson

import (
	"io"

	hjson "github.com/hjson/hjson-go/v4"
)

// NewEncoder constructs an HJSON encoder.
//
// This encoder is a thin adapter around `github.com/hjson/hjson-go` that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements HJSON encoding and decoding.
type Encoder struct{}

// Encode writes v to w as HJSON.
func (e *Encoder) Encode(w io.Writer, v any) error {
	data, err := hjson.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// Decode reads HJSON from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
func (e *Encoder) Decode(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return hjson.Unmarshal(data, v)
}
