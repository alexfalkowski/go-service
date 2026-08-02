package encoding

import (
	"maps"

	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/slices"
)

// NewMap constructs a Map with the default encoders.
//
// The returned registry includes these kinds, each registered under exactly one canonical name (no
// aliases, matching [github.com/alexfalkowski/go-service/v2/encoding/stream.Map]'s design):
// "json", "hjson", "yaml", "toml", "msgpack", "protobuf", "prototext", "protojson", "gob", "bytes".
//
// Callers that need to accept alternate spellings of these kinds (for example HTTP media subtypes such
// as "pb" or "octet-stream") are expected to translate them to the canonical kind above before calling
// [Map.Get]; see [github.com/alexfalkowski/go-service/v2/net/http/content/unary]'s unaryKind.
//
// Callers can add additional kinds or override existing kinds via [Map.Register].
func NewMap() *Map {
	return &Map{
		encoders: map[string]Encoder{
			"json":      json.NewEncoder(),
			"hjson":     hjson.NewEncoder(),
			"yaml":      yaml.NewEncoder(),
			"toml":      toml.NewEncoder(),
			"msgpack":   msgpack.NewEncoder(),
			"protobuf":  proto.NewBinary(),
			"prototext": proto.NewText(),
			"protojson": proto.NewJSON(),
			"gob":       gob.NewEncoder(),
			"bytes":     bytes.NewEncoder(),
		},
	}
}

// Map provides lookup and registration of encoders by kind.
//
// This type is a thin convenience around a string-keyed map and is commonly used with configuration
// to select an encoder at runtime.
//
// Map is not concurrency-safe. If you mutate it via Register, do so during initialization.
type Map struct {
	encoders map[string]Encoder
}

// Register associates kind with enc, overwriting any existing encoder.
//
// If kind already exists, the previous encoder is replaced.
func (f *Map) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get returns the encoder registered for kind.
//
// If no encoder is registered for kind, or if kind was registered with a nil encoder, Get returns nil.
// Callers typically treat nil as "unknown or unavailable kind" and fall back to a default encoder elsewhere.
func (f *Map) Get(kind string) Encoder {
	return f.encoders[kind]
}

// Keys returns the list of registered encoder kinds.
//
// Keys includes kinds registered with nil encoders. The returned slice is not guaranteed to be sorted.
func (f *Map) Keys() []string {
	return slices.Collect(maps.Keys(f.encoders))
}
