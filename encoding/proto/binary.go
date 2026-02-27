package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"google.golang.org/protobuf/proto"
)

// NewBinary constructs a protobuf binary encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/proto Marshal/Unmarshal that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewBinary() *Binary {
	return &Binary{}
}

// Binary implements protobuf binary encoding and decoding.
//
// Encode expects v to implement proto.Message and writes the protobuf binary wire format to the writer.
// Decode expects v to implement proto.Message (typically a pointer to a generated message) and unmarshals
// the protobuf binary wire format from the reader into v.
type Binary struct{}

// Encode writes v as protobuf binary (wire format) to w.
//
// If v does not implement proto.Message, Encode returns encoding/errors.ErrInvalidType.
// Any marshaling error from proto.Marshal and any write error from w.Write is returned.
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

// Decode reads protobuf binary (wire format) from r and unmarshals it into v.
//
// Decode reads all remaining bytes from r (via io.ReadAll) before unmarshaling.
// If v does not implement proto.Message, Decode returns encoding/errors.ErrInvalidType.
// Any read error from io.ReadAll and any unmarshal error from proto.Unmarshal is returned.
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
