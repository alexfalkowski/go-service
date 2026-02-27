package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewJSON constructs a protobuf JSON encoder.
//
// This encoder is a thin adapter around google.golang.org/protobuf/encoding/protojson Marshal/Unmarshal that
// satisfies `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
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
// If v does not implement proto.Message, Encode returns encoding/errors.ErrInvalidType.
// Any marshaling error from protojson.Marshal and any write error from w.Write is returned.
func (e *JSON) Encode(w io.Writer, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	bytes, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// Decode reads protobuf JSON from r and unmarshals it into v.
//
// Decode reads all remaining bytes from r (via io.ReadAll) before unmarshaling.
// If v does not implement proto.Message, Decode returns encoding/errors.ErrInvalidType.
// Any read error from io.ReadAll and any unmarshal error from protojson.Unmarshal is returned.
func (e *JSON) Decode(r io.Reader, v any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	msg, ok := v.(proto.Message)
	if !ok {
		return errors.ErrInvalidType
	}

	return protojson.Unmarshal(bytes, msg)
}
