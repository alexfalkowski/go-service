package crypto

import (
	"github.com/alexfalkowski/go-service/v3/crypto/aes"
	"github.com/alexfalkowski/go-service/v3/crypto/bcrypt"
	"github.com/alexfalkowski/go-service/v3/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v3/crypto/hmac"
	"github.com/alexfalkowski/go-service/v3/crypto/pem"
	"github.com/alexfalkowski/go-service/v3/crypto/rand"
	"github.com/alexfalkowski/go-service/v3/crypto/rsa"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires cryptographic subpackages into [github.com/alexfalkowski/go-signal/v2].
//
// It composes the modules from the crypto subpackages so services can include a single
// [Module] to register the supported crypto primitives and helpers.
//
// The included submodules currently wire:
//
//   - crypto/pem: PEM parsing helpers for keys/certificates.
//   - crypto/rand: cryptographically secure randomness utilities.
//   - crypto/aes: AES key generation and symmetric encryption helpers.
//   - crypto/bcrypt: password hashing helpers.
//   - crypto/ed25519: Ed25519 key generation, signing, and verification.
//   - crypto/hmac: HMAC key generation and signing/verification helpers.
//   - crypto/rsa: RSA key generation and encryption/decryption helpers.
//
// Note: this module only wires constructors; feature enablement is typically controlled via
// configuration in the consuming subsystems (for example nil/disabled sub-configs).
var Module = signal.Module(
	pem.Module,
	rand.Module,
	aes.Module,
	bcrypt.Module,
	ed25519.Module,
	hmac.Module,
	rsa.Module,
)
