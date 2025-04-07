package token

import (
	"github.com/alexfalkowski/go-service/token/ssh"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewKID),
	fx.Provide(NewJWT),
	fx.Provide(NewPaseto),
	fx.Provide(NewOpaque),
	ssh.Module,
	fx.Provide(NewToken),
)
