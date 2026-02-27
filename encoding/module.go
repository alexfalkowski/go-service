package encoding

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
)

// Module wires the default encoder implementations and the encoder registry into Fx/Dig.
//
// Provided constructors:
//
//   - Protobuf encoders:
//
//   - *proto.Binary (via proto.NewBinary)
//
//   - *proto.Text (via proto.NewText)
//
//   - *proto.JSON (via proto.NewJSON)
//
//   - Structured config encoders:
//
//   - *json.Encoder (via json.NewEncoder)
//
//   - *toml.Encoder (via toml.NewEncoder)
//
//   - *yaml.Encoder (via yaml.NewEncoder)
//
//   - Other encoders:
//
//   - *gob.Encoder (via gob.NewEncoder)
//
//   - *bytes.Encoder (via bytes.NewEncoder) for io.ReaderFrom/io.WriterTo passthrough
//
// Finally, it constructs an `*encoding.Map` via `NewMap`, pre-populated with common kind aliases
// (e.g. "yaml"/"yml", protobuf kind synonyms, and "plain"/"octet-stream" passthrough kinds).
var Module = di.Module(
	di.Constructor(proto.NewBinary),
	di.Constructor(proto.NewText),
	di.Constructor(proto.NewJSON),
	di.Constructor(json.NewEncoder),
	di.Constructor(toml.NewEncoder),
	di.Constructor(yaml.NewEncoder),
	di.Constructor(gob.NewEncoder),
	di.Constructor(bytes.NewEncoder),
	di.Constructor(NewMap),
)
