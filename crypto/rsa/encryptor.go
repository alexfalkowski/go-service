package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// NewEncryptor for rsa.
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

// Encryptor for rsa.
type Encryptor struct {
	generator *rand.Generator
	publicKey *rsa.PublicKey
}

// Encrypt for rsa.
func (e *Encryptor) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha512.New(), e.generator, e.publicKey, msg, nil)
}
