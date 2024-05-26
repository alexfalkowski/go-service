package security

import (
	"github.com/alexfalkowski/go-service/security/token"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	token.Module,
)
