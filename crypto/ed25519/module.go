package ed25519

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewSigner),
	di.Constructor(NewVerifier),
)
