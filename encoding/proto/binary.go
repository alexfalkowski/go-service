package proto

import (
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
	"google.golang.org/protobuf/proto"
)

// NewBinary constructs a protobuf binary encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/proto Marshal/Unmarshal that satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder].
func NewBinary() *Binary {
	return &Binary{}
}

// Binary implements protobuf binary encoding and decoding.
//
// Encode expects v to implement [proto.Message] and writes the protobuf binary wire format to the writer.
// Decode expects v to implement [proto.Message] (typically a pointer to a generated message) and unmarshals
// the protobuf binary wire format from the reader into v.
type Binary struct{}

// Encode writes v as protobuf binary (wire format) to w.
//
// If v does not implement [proto.Message], Encode returns [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType].
// A message with missing required fields fails to marshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
// Any marshaling error from [proto.MarshalOptions.Marshal] and any write error from w.Write is returned.
func (e *Binary) Encode(w io.Writer, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return marshalMessage(w, v, proto.MarshalOptions{AllowPartial: resolved.AllowPartial()}.Marshal)
}

// Decode reads protobuf binary (wire format) from r and unmarshals it into v.
//
// If v does not implement [proto.Message], Decode returns
// [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType] without reading from r.
//
// Decode otherwise reads all remaining bytes from r (via [io.ReadAll]) before
// unmarshaling. A message with missing required fields fails to unmarshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
//
// Any read error from [io.ReadAll] and any unmarshal error from [proto.UnmarshalOptions.Unmarshal] is returned.
func (e *Binary) Decode(r io.Reader, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return unmarshalMessage(r, v, proto.UnmarshalOptions{AllowPartial: resolved.AllowPartial()}.Unmarshal)
}
