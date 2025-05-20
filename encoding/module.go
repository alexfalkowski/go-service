package encoding

import (
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(proto.NewBinary),
	fx.Provide(proto.NewText),
	fx.Provide(proto.NewJSON),
	fx.Provide(json.NewEncoder),
	fx.Provide(toml.NewEncoder),
	fx.Provide(yaml.NewEncoder),
	fx.Provide(gob.NewEncoder),
	fx.Provide(bytes.NewEncoder),
	fx.Provide(NewMap),
)
