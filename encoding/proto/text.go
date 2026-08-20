package proto

import (
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
	"google.golang.org/protobuf/encoding/prototext"
)

// NewText constructs a protobuf text encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/encoding/prototext Marshal/Unmarshal that
// satisfies [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder].
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
// A message with missing required fields fails to marshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
// Any marshaling error from [prototext.MarshalOptions.Marshal] and any write error from w.Write is returned.
func (e *Text) Encode(w io.Writer, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return marshalMessage(w, v, prototext.MarshalOptions{AllowPartial: resolved.AllowPartial()}.Marshal)
}

// Decode reads protobuf text format from r and unmarshals it into v.
//
// If v does not implement proto.Message, Decode returns
// [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType] without reading from r.
//
// Decode otherwise reads all remaining bytes from r (via [io.ReadAll]) before
// unmarshaling. Unknown protobuf text fields are rejected unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] is
// supplied. A message with missing required fields fails to unmarshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
//
// Any read error from [io.ReadAll] and any unmarshal error from [prototext.UnmarshalOptions.Unmarshal] is returned.
func (e *Text) Decode(r io.Reader, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return unmarshalMessage(r, v, prototext.UnmarshalOptions{DiscardUnknown: resolved.DiscardUnknown(), AllowPartial: resolved.AllowPartial()}.Unmarshal)
}
