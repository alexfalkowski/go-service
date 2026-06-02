package encoding

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
)

// Module wires the default encoder implementations and the encoder registry into [go.uber.org/fx]/[go.uber.org/dig].
//
// It provides protobuf encoders (*[proto.Binary], *[proto.Text], *[proto.JSON]),
// structured config encoders (*[json.Encoder], *[hjson.Encoder], *[toml.Encoder],
// *[yaml.Encoder], *[msgpack.Encoder]), and passthrough/specialized encoders
// (*[gob.Encoder] and *[bytes.Encoder]).
//
// Finally, it constructs a *[Map] via [NewMap], pre-populated with common kind
// aliases (for example "yaml"/"yml", protobuf kind synonyms, and
// "plain"/"octet-stream" passthrough kinds).
var Module = di.Module(
	di.Constructor(proto.NewBinary),
	di.Constructor(proto.NewText),
	di.Constructor(proto.NewJSON),
	di.Constructor(json.NewEncoder),
	di.Constructor(hjson.NewEncoder),
	di.Constructor(toml.NewEncoder),
	di.Constructor(yaml.NewEncoder),
	di.Constructor(msgpack.NewEncoder),
	di.Constructor(gob.NewEncoder),
	di.Constructor(bytes.NewEncoder),
	di.Constructor(NewMap),
)
