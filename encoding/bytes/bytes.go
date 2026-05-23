package bytes

import (
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/reflect"
)

// NewEncoder constructs a passthrough encoder for stream-capable types.
//
// This encoder is intended for values that can write themselves to an io.Writer
// and/or read themselves from an io.Reader via the shared go-service io package
// aliases:
//
//   - io.WriterTo for Encode
//   - io.ReaderFrom for Decode
//
// It is useful when you want to treat a value as its own codec (for example when caching or transporting
// pre-serialized payloads) while still satisfying the go-service encoding.Encoder interface.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder encodes and decodes values by delegating to io.WriterTo and io.ReaderFrom.
//
// This encoder does not perform any framing, escaping, or format conversion. It simply forwards the
// Encode/Decode call to the value itself.
type Encoder struct{}

// Encode writes v to w when v implements io.WriterTo.
//
// If v does not implement io.WriterTo, Encode returns encoding/errors.ErrInvalidType.
// Any error returned by WriteTo is returned.
func (e *Encoder) Encode(w io.Writer, v any) error {
	to, err := writerTo(v)
	if err != nil {
		return err
	}

	_, err = to.WriteTo(w)
	return err
}

// Decode reads from r into v when v implements io.ReaderFrom.
//
// If v does not implement io.ReaderFrom, Decode returns encoding/errors.ErrInvalidType.
// If v also implements io.Resetter, Decode resets it before reading so the
// destination is repopulated rather than appended to.
// Any error returned by ReadFrom is returned.
func (e *Encoder) Decode(r io.Reader, v any) error {
	from, err := readerFrom(v)
	if err != nil {
		return err
	}

	if reset, ok := v.(io.Resetter); ok && !reflect.IsNil(reset) {
		reset.Reset()
	}

	_, err = from.ReadFrom(r)
	return err
}

func writerTo(v any) (io.WriterTo, error) {
	to, ok := v.(io.WriterTo)
	if !ok || reflect.IsNil(to) {
		return nil, errors.ErrInvalidType
	}

	return to, nil
}

func readerFrom(v any) (io.ReaderFrom, error) {
	from, ok := v.(io.ReaderFrom)
	if !ok || reflect.IsNil(from) {
		return nil, errors.ErrInvalidType
	}

	return from, nil
}
