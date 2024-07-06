package proto

import (
	"google.golang.org/protobuf/proto"
)

// Marshaller for proto.
type Marshaller struct{}

// NewMarshaller for proto.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
