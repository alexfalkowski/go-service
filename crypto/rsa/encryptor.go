package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// NewEncryptor constructs an Encryptor when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, it loads the public key using cfg.PublicKey.
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

// Encrypt encrypts msg using RSA-OAEP with SHA-512.
//
// The OAEP label is nil.
func (e *Encryptor) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), e.generator, e.publicKey, msg, nil)
}
