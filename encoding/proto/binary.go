package proto

import (
	"io"

	"github.com/alexfalkowski/go-service/errors"
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
			err = errors.Prefix("proto", runtime.ConvertRecover(r))
		}
	}()

	b, err := proto.Marshal(v.(proto.Message))
	runtime.Must(err)

	_, err = w.Write(b)
	runtime.Must(err)

	return
}

func (e *Binary) Decode(r io.Reader, v any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("proto", runtime.ConvertRecover(r))
		}
	}()

	b, err := io.ReadAll(r)
	runtime.Must(err)

	err = proto.Unmarshal(b, v.(proto.Message))
	runtime.Must(err)

	return
}
