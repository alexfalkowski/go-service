package bytes

import (
	"io"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/encoding/errors"
)

// NewEncoder for bytes.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder for bytes.
type Encoder struct{}

// Encode for bytes.
func (e *Encoder) Encode(w io.Writer, v any) error {
	buffer, ok := v.(*bytes.Buffer)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := io.Copy(w, buffer)

	return err
}

// Decode for bytes.
func (e *Encoder) Decode(r io.Reader, v any) error {
	buffer, ok := v.(*bytes.Buffer)
	if !ok {
		return errors.ErrInvalidType
	}

	_, err := io.Copy(buffer, r)

	return err
}
