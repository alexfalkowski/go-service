package token

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/token/jwt"
	"github.com/alexfalkowski/go-service/v2/token/paseto"
	"github.com/alexfalkowski/go-service/v2/token/ssh"
)

// Module for fx.
var Module = di.Module(
	access.Module,
	jwt.Module,
	paseto.Module,
	ssh.Module,
	di.Constructor(NewToken),
	di.Constructor(NewGenerator),
	di.Constructor(NewVerifier),
)
