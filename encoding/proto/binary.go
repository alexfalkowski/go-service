package proto

import (
	"io"

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
	b, err := proto.Marshal(v.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

// Decode for proto.
func (e *Binary) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return proto.Unmarshal(b, v.(proto.Message))
}
