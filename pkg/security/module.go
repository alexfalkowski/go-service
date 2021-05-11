package security

import (
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"go.uber.org/fx"
)

// Auth0Module for fx.
var Auth0Module = fx.Options(fx.Provide(auth0.NewConfig))
