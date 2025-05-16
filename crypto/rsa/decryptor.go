package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/crypto/pem"
	"github.com/alexfalkowski/go-service/crypto/rand"
)

// NewDecryptor for rsa.
func NewDecryptor(generator *rand.Generator, decoder *pem.Decoder, cfg *Config) (*Decryptor, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	pri, err := cfg.PrivateKey(decoder)
	if err != nil {
		return nil, err
	}

	return &Decryptor{generator: generator, privateKey: pri}, nil
}

// Cipher for rsa.
type Decryptor struct {
	generator  *rand.Generator
	privateKey *rsa.PrivateKey
}

// Decrypt for rsa.
func (d *Decryptor) Decrypt(msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), d.generator, d.privateKey, msg, nil)
}
