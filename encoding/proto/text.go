package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/runtime"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// NewText for proto.
func NewText() *Text {
	return &Text{}
}

// Text for proto.
type Text struct{}

// Encode for proto.
func (e *Text) Encode(w io.Writer, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	b, err := prototext.Marshal(v.(proto.Message))
	runtime.Must(err)

	_, err = w.Write(b)
	runtime.Must(err)

	return
}

// Decode for proto.
func (e *Text) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return prototext.Unmarshal(b, v.(proto.Message))
}
