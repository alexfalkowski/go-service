package token

import (
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	jwt.Module,
	paseto.Module,
	ssh.Module,
	fx.Provide(NewToken),
)
