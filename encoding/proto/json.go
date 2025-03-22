package proto

import (
	"io"

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
	b, err := protojson.Marshal(v.(proto.Message))
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

// Decode for proto.
func (e *JSON) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return protojson.Unmarshal(b, v.(proto.Message))
}
