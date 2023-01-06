package marshaller

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPack for marshaller.
type MsgPack struct{}

// NewMsgPack for marshaller.
func NewMsgPack() *MsgPack {
	return &MsgPack{}
}

func (m *MsgPack) Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *MsgPack) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
