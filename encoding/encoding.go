package encoding

import (
	"maps"
	"slices"

	"github.com/alexfalkowski/go-service/encoding/bytes"
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
	Bytes       *bytes.Encoder
}

// NewMap for encoding.
func NewMap(params MapParams) *Map {
	return &Map{
		encoders: map[string]Encoder{
			"json":         params.JSON,
			"yaml":         params.YAML,
			"yml":          params.YAML,
			"toml":         params.TOML,
			"pb":           params.ProtoBinary,
			"pbbin":        params.ProtoBinary,
			"proto":        params.ProtoBinary,
			"protobin":     params.ProtoBinary,
			"protobuf":     params.ProtoBinary,
			"pbtxt":        params.ProtoText,
			"prototext":    params.ProtoText,
			"prototxt":     params.ProtoText,
			"protojson":    params.ProtoJSON,
			"pbjson":       params.ProtoJSON,
			"gob":          params.GOB,
			"markdown":     params.Bytes,
			"octet-stream": params.Bytes,
			"plain":        params.Bytes,
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

// Keys from the encoders map.
func (f *Map) Keys() []string {
	return slices.Collect(maps.Keys(f.encoders))
}
