package rsa

import (
	"crypto/rsa"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/message"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/strings"
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
	if strings.IsEmpty(cfg.Public) {
		return nil, errors.ErrMissingKey
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

// Encrypt encrypts msg.Data using RSA-OAEP with SHA-512 and returns the ciphertext.
//
// OAEP parameters:
//   - hash: SHA-512
//   - label: msg.Meta
//
// msg.Meta is authenticated context and must be supplied unchanged to Decrypt.
// msg.Data must fit RSA-OAEP's plaintext limit for the public key size:
// modulus bytes minus two SHA-512 digest lengths minus two bytes. For this
// package's default 4096-bit keys, the maximum plaintext size is 382 bytes.
//
// The randomness source is the injected generator's reader.
func (e *Encryptor) Encrypt(msg message.Message) ([]byte, error) {
	return EncryptOAEP(e.generator, e.publicKey, msg)
}
