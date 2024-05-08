package crypto

import (
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	aes.Module,
	argon2.Module,
	ed25519.Module,
	hmac.Module,
	rsa.Module,
)
