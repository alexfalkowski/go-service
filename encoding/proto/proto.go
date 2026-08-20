package proto

import (
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/reflect"
	"google.golang.org/protobuf/proto"
)

// Message is an alias for proto.Message.
type Message = proto.Message

func message(v any) (Message, error) {
	msg, ok := v.(Message)
	if !ok || reflect.IsNil(msg) {
		return nil, errors.ErrInvalidType
	}

	return msg, nil
}

func readMessage(r io.Reader, v any) (Message, []byte, error) {
	msg, err := message(v)
	if err != nil {
		return nil, nil, err
	}

	bytes, _, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	return msg, bytes, nil
}
