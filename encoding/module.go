package encoding

import (
	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
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
	fx.Provide(NewMap),
)
