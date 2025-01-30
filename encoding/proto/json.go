package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// NewJSON for proto.
func NewJSON() *JSON {
	return &JSON{}
}

// JSON for proto.
type JSON struct{}

func (e *JSON) Encode(w io.Writer, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	b, err := protojson.Marshal(v.(proto.Message))
	runtime.Must(err)

	_, err = w.Write(b)
	runtime.Must(err)

	return
}

func (e *JSON) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return protojson.Unmarshal(b, v.(proto.Message))
}
