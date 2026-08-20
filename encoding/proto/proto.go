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

func marshalMessage(w io.Writer, v any, marshal func(Message) ([]byte, error)) error {
	msg, err := message(v)
	if err != nil {
		return err
	}

	bytes, err := marshal(msg)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

func unmarshalMessage(r io.Reader, v any, unmarshal func([]byte, Message) error) error {
	msg, bytes, err := readMessage(r, v)
	if err != nil {
		return err
	}

	return unmarshal(bytes, msg)
}
