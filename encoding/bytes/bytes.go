package bytes

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
)

// NewEncoder for bytes.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for bytes.
type Encoder struct{}

// Encode for bytes.
func (e *Encoder) Encode(w io.Writer, v any) error {
	to, ok := v.(io.WriterTo)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := to.WriteTo(w)
	return err
}

// Decode for bytes.
func (e *Encoder) Decode(r io.Reader, v any) error {
	from, ok := v.(io.ReaderFrom)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := from.ReadFrom(r)
	return err
}
