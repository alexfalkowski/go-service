package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewJSON constructs a protobuf JSON encoder.
func NewJSON() *JSON {
	return &JSON{}
}

// JSON implements protobuf JSON encoding and decoding.
type JSON struct{}

// Encode writes v as protobuf JSON to w.
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

// Decode reads protobuf JSON from r into v.
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
