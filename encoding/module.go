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
	fx.Provide(proto.NewMarshaller),
	fx.Provide(json.NewMarshaller),
	fx.Provide(toml.NewMarshaller),
	fx.Provide(yaml.NewMarshaller),
	fx.Provide(gob.NewMarshaller),
	fx.Provide(NewMap),
)
