package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/runtime"
	"google.golang.org/protobuf/proto"
)

// NewBinary for proto.
func NewBinary() *Binary {
	return &Binary{}
}

// Binary for proto.
type Binary struct{}

func (e *Binary) Encode(w io.Writer, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = runtime.ConvertRecover(r)
		}
	}()

	b, err := proto.Marshal(v.(proto.Message))
	runtime.Must(err)

	_, err = w.Write(b)
	runtime.Must(err)

	return
}

func (e *Binary) Decode(r io.Reader, v any) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return proto.Unmarshal(b, v.(proto.Message))
}
