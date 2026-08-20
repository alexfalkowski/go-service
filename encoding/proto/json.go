package proto

import (
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewJSON constructs a protobuf JSON encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/encoding/protojson Marshal/Unmarshal that
// satisfies [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder].
func NewJSON() *JSON {
	return &JSON{}
}

// JSON implements protobuf JSON encoding and decoding.
//
// Encode expects v to implement proto.Message and writes protobuf JSON to the writer.
// Decode expects v to implement proto.Message (typically a pointer to a generated message) and unmarshals
// protobuf JSON from the reader into v.
type JSON struct{}

// Encode writes v as protobuf JSON to w.
//
// If v does not implement proto.Message, Encode returns [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType].
// A message with missing required fields fails to marshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
// Any marshaling error from [protojson.MarshalOptions.Marshal] and any write error from w.Write is returned.
func (e *JSON) Encode(w io.Writer, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return marshalMessage(w, v, protojson.MarshalOptions{AllowPartial: resolved.AllowPartial()}.Marshal)
}

// Decode reads protobuf JSON from r and unmarshals it into v.
//
// If v does not implement proto.Message, Decode returns
// [github.com/alexfalkowski/go-service/v2/encoding/errors.ErrInvalidType] without reading from r.
//
// Decode otherwise reads all remaining bytes from r (via [io.ReadAll]) before
// unmarshaling. Unknown protobuf JSON fields are rejected unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] is
// supplied. A message with missing required fields fails to unmarshal unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] is supplied.
//
// Any read error from [io.ReadAll] and any unmarshal error from [protojson.UnmarshalOptions.Unmarshal] is returned.
func (e *JSON) Decode(r io.Reader, v any, opts ...codec.Option) error {
	resolved := codec.Apply(opts...)

	return unmarshalMessage(r, v, protojson.UnmarshalOptions{DiscardUnknown: resolved.DiscardUnknown(), AllowPartial: resolved.AllowPartial()}.Unmarshal)
}
