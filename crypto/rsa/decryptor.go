package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// NewDecryptor constructs a Decryptor when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, it loads the private key using cfg.PrivateKey.
func NewDecryptor(generator *rand.Generator, decoder *pem.Decoder, cfg *Config) (*Decryptor, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pri, err := cfg.PrivateKey(decoder)
	if err != nil {
		return nil, err
	}

	return &Decryptor{generator: generator, privateKey: pri}, nil
}

// Decryptor holds an RSA private key and randomness source used for decryption.
type Decryptor struct {
	generator  *rand.Generator
	privateKey *rsa.PrivateKey
}

// Decrypt decrypts msg using RSA-OAEP with SHA-512.
//
// The OAEP label is nil.
func (d *Decryptor) Decrypt(msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), d.generator, d.privateKey, msg, nil)
}
