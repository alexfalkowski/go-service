package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// NewEncryptor constructs an RSA Encryptor when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewEncryptor returns (nil, nil).
//
// Enabled behavior: NewEncryptor loads and parses the RSA public key via cfg.PublicKey(decoder) and returns
// an Encryptor that performs RSA-OAEP encryption.
//
// Any error encountered while decoding/parsing the public key is returned.
func NewEncryptor(generator *rand.Generator, decoder *pem.Decoder, cfg *Config) (*Encryptor, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pub, err := cfg.PublicKey(decoder)
	if err != nil {
		return nil, err
	}

	return &Encryptor{generator: generator, publicKey: pub}, nil
}

// Encryptor holds an RSA public key and randomness source used for encryption.
type Encryptor struct {
	generator *rand.Generator
	publicKey *rsa.PublicKey
}

// Encrypt encrypts msg using RSA-OAEP with SHA-512 and returns the ciphertext.
//
// OAEP parameters:
//   - hash: SHA-512
//   - label: nil
//
// The randomness source is the injected crypto/rand generator.
func (e *Encryptor) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), e.generator, e.publicKey, msg, nil)
}
