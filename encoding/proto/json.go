package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/encoding/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewJSON for proto.
func NewJSON() *JSON {
	return &JSON{}
}

// JSON for proto.
type JSON struct{}

// Encode for proto.
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

// Decode for proto.
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
