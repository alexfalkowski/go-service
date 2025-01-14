package crypto

import (
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/argon2"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	rand.Module,
	aes.Module,
	argon2.Module,
	ed25519.Module,
	hmac.Module,
	rsa.Module,
	ssh.Module,
)
