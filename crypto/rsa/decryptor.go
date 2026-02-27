package rsa

import (
	"crypto/rsa"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
)

// NewDecryptor constructs an RSA Decryptor when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewDecryptor returns (nil, nil).
//
// Enabled behavior: NewDecryptor loads and parses the RSA private key via cfg.PrivateKey(decoder) and returns
// a Decryptor that performs RSA-OAEP decryption.
//
// Any error encountered while decoding/parsing the private key is returned.
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

// Decrypt decrypts msg using RSA-OAEP with SHA-512 and returns the plaintext.
//
// OAEP parameters:
//   - hash: SHA-512
//   - label: nil
//
// The randomness source is the injected crypto/rand generator.
func (d *Decryptor) Decrypt(msg []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha512.New(), d.generator, d.privateKey, msg, nil)
}
