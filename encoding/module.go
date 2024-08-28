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
	fx.Provide(proto.NewEncoder),
	fx.Provide(json.NewMarshaller),
	fx.Provide(json.NewEncoder),
	fx.Provide(toml.NewMarshaller),
	fx.Provide(toml.NewEncoder),
	fx.Provide(yaml.NewMarshaller),
	fx.Provide(yaml.NewEncoder),
	fx.Provide(gob.NewMarshaller),
	fx.Provide(gob.NewEncoder),
	fx.Provide(NewMarshallerMap),
	fx.Provide(NewEncoderMap),
)
