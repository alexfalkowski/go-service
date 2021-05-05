package security

import (
	"github.com/alexfalkowski/go-service/pkg/security/token"
	"go.uber.org/fx"
)

var (
	// NoopTokenGeneratorModule for fx.
	NoopTokenGeneratorModule = fx.Provide(token.NewNoopGenerator)

	// NoopTokenVerifierModule for fx.
	NoopTokenVerifierModule = fx.Provide(token.NewNoopVerifier)

	// NoopTokenModule for fx.
	NoopTokenModule = fx.Options(NoopTokenVerifierModule, NoopTokenGeneratorModule)
)
