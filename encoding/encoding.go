package encoding

import (
	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
	"go.uber.org/fx"
)

// MapParams for encoding.
type MapParams struct {
	fx.In

	JSON        *json.Encoder
	YAML        *yaml.Encoder
	TOML        *toml.Encoder
	ProtoBinary *proto.Binary
	ProtoText   *proto.Text
	ProtoJSON   *proto.JSON
	GOB         *gob.Encoder
}

// NewMap for encoding.
func NewMap(params MapParams) *Map {
	return &Map{
		encoders: map[string]Encoder{
			"json":      params.JSON,
			"yaml":      params.YAML,
			"yml":       params.YAML,
			"toml":      params.TOML,
			"proto":     params.ProtoBinary,
			"protobuf":  params.ProtoBinary,
			"prototext": params.ProtoText,
			"protojson": params.ProtoJSON,
			"gob":       params.GOB,
		},
	}
}

// Map of encoding.
type Map struct {
	encoders map[string]Encoder
}

// Register kind and encoder.
func (f *Map) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get from kind.
func (f *Map) Get(kind string) Encoder {
	return f.encoders[kind]
}
