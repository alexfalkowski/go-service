package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// NewText for proto.
func NewText() *Text {
	return &Text{}
}

// Text for proto.
type Text struct{}

// Encode for proto.
func (e *Text) Encode(w io.Writer, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	bytes, err := prototext.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)

	return err
}

// Decode for proto.
func (e *Text) Decode(r io.Reader, v any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	return prototext.Unmarshal(bytes, msg)
}
