package token

import (
	"github.com/alexfalkowski/go-service/token/jwt"
	"github.com/alexfalkowski/go-service/token/opaque"
	"github.com/alexfalkowski/go-service/token/paseto"
	"github.com/alexfalkowski/go-service/token/ssh"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	jwt.Module,
	opaque.Module,
	paseto.Module,
	ssh.Module,
	fx.Provide(NewToken),
)
