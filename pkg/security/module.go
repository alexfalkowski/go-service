package security

import (
	"github.com/alexfalkowski/go-service/pkg/security/auth0"
	"go.uber.org/fx"
)

// Auth0Module for fx.
// nolint:gochecknoglobals
var Auth0Module = fx.Options(fx.Provide(auth0.NewGenerator), fx.Provide(auth0.NewCertificator), fx.Provide(auth0.NewVerifier))
