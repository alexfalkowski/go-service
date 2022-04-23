package marshaller

import (
	"github.com/vmihailenco/msgpack/v5"
)

type msgpackMarshaller struct{}

// NewMsgPack for marshaller.
// nolint:ireturn
func NewMsgPack() Marshaller {
	return &msgpackMarshaller{}
}

func (m *msgpackMarshaller) Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *msgpackMarshaller) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
