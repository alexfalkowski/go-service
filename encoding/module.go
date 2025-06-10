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

// Module for fx.
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
