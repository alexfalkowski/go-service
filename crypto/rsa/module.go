package rsa

import "github.com/alexfalkowski/go-service/v2/di"

// Module for fx.
var Module = di.Module(
	di.Constructor(NewGenerator),
	di.Constructor(NewEncryptor),
	di.Constructor(NewDecryptor),
)
