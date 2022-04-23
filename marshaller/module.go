package marshaller

import (
	"go.uber.org/fx"
)

// ProtoModule for fx.
// nolint:gochecknoglobals
var ProtoModule = fx.Options(fx.Provide(NewProto))
