package token

import (
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewKID),
	fx.Provide(NewJWT),
	fx.Provide(NewPaseto),
	fx.Provide(NewToken),
)
