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

// Module wires cryptographic subpackages into Fx/Dig.
//
// It composes the modules from the crypto subpackages so services can include a single
// `crypto.Module` to register the supported crypto primitives and helpers.
//
// The included submodules currently wire:
//
//   - `crypto/pem`: PEM parsing helpers for keys/certificates.
//   - `crypto/rand`: cryptographically secure randomness utilities.
//   - `crypto/aes`: AES key generation and symmetric encryption helpers.
//   - `crypto/bcrypt`: password hashing helpers.
//   - `crypto/ed25519`: Ed25519 key generation, signing, and verification.
//   - `crypto/hmac`: HMAC key generation and signing/verification helpers.
//   - `crypto/rsa`: RSA key generation and encryption/decryption helpers.
//   - `crypto/ssh`: SSH key generation, signing, and verification.
//
// Note: this module only wires constructors; feature enablement is typically controlled via
// configuration in the consuming subsystems (e.g. nil/disabled sub-configs).
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
