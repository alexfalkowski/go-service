package proto

import (
	"io"

	"google.golang.org/protobuf/proto"
)

// Encoder for proto.
type Encoder struct{}

// NewEncoder for proto.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	b, err := proto.Marshal(v.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return proto.Unmarshal(b, v.(proto.Message))
}
