package marshaller

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewMsgPack),
	fx.Provide(NewProto),
	fx.Provide(NewYAML),
	fx.Provide(NewTOML),
	fx.Provide(NewFactory),
)
