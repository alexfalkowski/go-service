package rsa

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewGenerator),
	fx.Provide(NewEncryptor),
	fx.Provide(NewDecryptor),
)
