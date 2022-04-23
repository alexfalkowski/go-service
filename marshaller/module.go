package marshaller

import (
	"go.uber.org/fx"
)

var (
	// ProtoModule for fx.
	// nolint:gochecknoglobals
	ProtoModule = fx.Options(fx.Provide(NewProto))

	// MsgPackModule for fx.
	// nolint:gochecknoglobals
	MsgPackModule = fx.Options(fx.Provide(NewMsgPack))
)
