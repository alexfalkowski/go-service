package bytes

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
)

// NewEncoder constructs a bytes encoder for io.ReaderFrom/io.WriterTo types.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder encodes and decodes using io.ReaderFrom/io.WriterTo.
type Encoder struct{}

// Encode writes v to w when v implements io.WriterTo.
func (e *Encoder) Encode(w io.Writer, v any) error {
	to, ok := v.(io.WriterTo)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := to.WriteTo(w)
	return err
}

// Decode reads from r into v when v implements io.ReaderFrom.
func (e *Encoder) Decode(r io.Reader, v any) error {
	from, ok := v.(io.ReaderFrom)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := from.ReadFrom(r)
	return err
}
