package crypto

import (
	"github.com/alexfalkowski/go-service/crypto/aes"
	"github.com/alexfalkowski/go-service/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/hmac"
	"github.com/alexfalkowski/go-service/crypto/pem"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/crypto/rsa"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	pem.Module,
	rand.Module,
	aes.Module,
	bcrypt.Module,
	ed25519.Module,
	hmac.Module,
	rsa.Module,
	ssh.Module,
)
