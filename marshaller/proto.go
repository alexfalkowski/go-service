package marshaller

import (
	"google.golang.org/protobuf/proto"
)

type protoMarshaller struct{}

// NewProto for marshaller.
func NewProto() Marshaller {
	return &protoMarshaller{}
}

func (m *protoMarshaller) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (m *protoMarshaller) Unmarshal(data []byte, v any) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
