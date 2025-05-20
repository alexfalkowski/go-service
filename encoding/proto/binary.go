package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"google.golang.org/protobuf/proto"
)

// NewBinary for proto.
func NewBinary() *Binary {
	return &Binary{}
}

// Binary for proto.
type Binary struct{}

// Encode for proto.
func (e *Binary) Encode(w io.Writer, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)

	return err
}

// Decode for proto.
func (e *Binary) Decode(r io.Reader, v any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	return proto.Unmarshal(bytes, msg)
}
