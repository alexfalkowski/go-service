package crypto

import (
	"github.com/alexfalkowski/go-service/v2/crypto/aes"
	"github.com/alexfalkowski/go-service/v2/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/crypto/hmac"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/crypto/rsa"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module for fx.
var Module = di.Module(
	pem.Module,
	rand.Module,
	aes.Module,
	bcrypt.Module,
	ed25519.Module,
	hmac.Module,
	rsa.Module,
	ssh.Module,
)
