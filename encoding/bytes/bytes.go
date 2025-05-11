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
	switch kind := v.(type) {
	case *[]byte:
		_, err := w.Write(*kind)

		return err
	case *bytes.Buffer:
		_, err := io.Copy(w, kind)

		return err
	default:
		return errors.ErrInvalidType
	}
}

// Decode for bytes.
func (e *Encoder) Decode(r io.Reader, v any) error {
	switch kind := v.(type) {
	case *[]byte:
		_, err := r.Read(*kind)

		return err
	case *bytes.Buffer:
		_, err := io.Copy(kind, r)

		return err
	default:
		return errors.ErrInvalidType
	}
}
