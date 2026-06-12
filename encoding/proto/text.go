package proto

import (
	"github.com/alexfalkowski/go-service/v2/io"
	"google.golang.org/protobuf/encoding/prototext"
)

// NewText constructs a protobuf text encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/encoding/prototext Marshal/Unmarshal that
// satisfies [github.com/alexfalkowski/go-service/v2/encoding.Encoder].
func NewText() *Text {
	return &Text{}
}

// Text implements protobuf text encoding and decoding.
//
// Encode expects v to implement proto.Message and writes protobuf text format to the writer.
// Decode expects v to implement proto.Message (typically a pointer to a generated message) and unmarshals
// protobuf text format from the reader into v.
type Text struct{}

// Encode writes v as protobuf text format to w.
//
// If v does not implement proto.Message, Encode returns [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType].
// Any marshaling error from [prototext.Marshal] and any write error from w.Write is returned.
func (e *Text) Encode(w io.Writer, v any) error {
	msg, err := message(v)
	if err != nil {
		return err
	}

	bytes, err := prototext.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// Decode reads protobuf text format from r and unmarshals it into v.
//
// If v does not implement proto.Message, Decode returns
// [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType] without reading from r.
//
// Decode otherwise reads all remaining bytes from r (via [io.ReadAll]) before
// unmarshaling.
//
// Any read error from [io.ReadAll] and any unmarshal error from [prototext.UnmarshalOptions.Unmarshal] is returned.
func (e *Text) Decode(r io.Reader, v any) error {
	msg, err := message(v)
	if err != nil {
		return err
	}

	bytes, _, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return prototext.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(bytes, msg)
}
