package proto

import (
	"io"

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
	b, err := prototext.Marshal(v.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

// Decode for proto.
func (e *Text) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return prototext.Unmarshal(b, v.(proto.Message))
}
