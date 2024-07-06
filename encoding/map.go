package encoding

import (
	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
)

type configs map[string]Marshaller

// Map of marshaller.
type Map struct {
	configs configs
}

// NewMap for marshaller.
func NewMap() *Map {
	f := &Map{
		configs: configs{
			"json":     json.NewMarshaller(),
			"yaml":     yaml.NewMarshaller(),
			"yml":      yaml.NewMarshaller(),
			"toml":     toml.NewMarshaller(),
			"proto":    proto.NewMarshaller(),
			"protobuf": proto.NewMarshaller(),
			"gob":      gob.NewMarshaller(),
		},
	}

	return f
}

// Register kind and marshaller.
func (f *Map) Register(kind string, m Marshaller) {
	f.configs[kind] = m
}

// Get from kind.
func (f *Map) Get(kind string) Marshaller {
	return f.configs[kind]
}
