package proto

import (
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/reflect"
	"google.golang.org/protobuf/proto"
)

func message(v any) (proto.Message, error) {
	msg, ok := v.(proto.Message)
	if !ok || reflect.IsNil(msg) {
		return nil, errors.ErrInvalidType
	}

	return msg, nil
}
