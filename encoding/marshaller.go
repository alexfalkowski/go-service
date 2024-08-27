package encoding

import (
	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
)

// Marshaller allows to have different ways to marshal/unmarshal.
type Marshaller interface {
	// Marshal value.
	Marshal(v any) ([]byte, error)

	// Unmarshal data to value.
	Unmarshal(data []byte, v any) error
}

type marshallers map[string]Marshaller

// MarshallerMap of marshaller.
type MarshallerMap struct {
	marshallers marshallers
}

// NewMarshallerMap for marshaller.
func NewMarshallerMap() *MarshallerMap {
	m := &MarshallerMap{
		marshallers: marshallers{
			"json":     json.NewMarshaller(),
			"yaml":     yaml.NewMarshaller(),
			"yml":      yaml.NewMarshaller(),
			"toml":     toml.NewMarshaller(),
			"proto":    proto.NewMarshaller(),
			"protobuf": proto.NewMarshaller(),
			"gob":      gob.NewMarshaller(),
		},
	}

	return m
}

// Register kind and marshaller.
func (m *MarshallerMap) Register(kind string, mar Marshaller) {
	m.marshallers[kind] = mar
}

// Get from kind.
func (m *MarshallerMap) Get(kind string) Marshaller {
	return m.marshallers[kind]
}
