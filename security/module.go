package security

import (
	"github.com/alexfalkowski/go-service/security/auth0"
	"go.uber.org/fx"
)

// Auth0Module for fx.
var Auth0Module = fx.Options(fx.Provide(auth0.NewGenerator), fx.Provide(auth0.NewCertificator), fx.Provide(auth0.NewVerifier))
