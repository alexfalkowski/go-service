package marshaller

import (
	"google.golang.org/protobuf/proto"
)

type protoMarshaller struct{}

// NewProto for marshaller.
// nolint:ireturn
func NewProto() Marshaller {
	return &protoMarshaller{}
}

// nolint:forcetypeassert
func (m *protoMarshaller) Marshal(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

// nolint:forcetypeassert
func (m *protoMarshaller) Unmarshal(data []byte, v any) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
