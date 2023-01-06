package marshaller

import (
	"google.golang.org/protobuf/proto"
)

// Proto for marshaller.
type Proto struct{}

// NewProto for marshaller.
func NewProto() *Proto {
	return &Proto{}
}

func (m *Proto) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (m *Proto) Unmarshal(data []byte, v any) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
