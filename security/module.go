package security

import (
	"github.com/alexfalkowski/go-service/security/oauth"
	"go.uber.org/fx"
)

// OAuthModule for fx.
var OAuthModule = fx.Options(
	fx.Provide(oauth.NewGenerator),
	fx.Provide(oauth.NewCertificator),
	fx.Provide(oauth.NewVerifier),
)
