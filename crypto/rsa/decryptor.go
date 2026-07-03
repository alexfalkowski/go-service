package rsa

import (
	"crypto/rsa"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/message"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/strings"
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
	if strings.IsEmpty(cfg.Private) {
		return nil, errors.ErrMissingKey
	}

	pri, err := cfg.PrivateKey(decoder)
	if err != nil {
		return nil, err
	}

	return &Decryptor{generator: generator, privateKey: pri}, nil
}

// Decryptor holds an RSA private key used for decryption.
//
// The generator is retained for API symmetry with Encryptor. The standard library's OAEP
// decryption treats its random parameter as legacy and ignores it.
type Decryptor struct {
	generator  *rand.Generator
	privateKey *rsa.PrivateKey
}

// Decrypt decrypts msg.Data using RSA-OAEP with SHA-512 and returns the plaintext.
//
// OAEP parameters:
//   - hash: SHA-512
//   - label: msg.Meta
//
// msg.Meta must match the metadata supplied to Encrypt.
// The injected generator is passed through for API consistency, but RSA-OAEP decryption ignores
// the randomness parameter.
func (d *Decryptor) Decrypt(msg message.Message) ([]byte, error) {
	return DecryptOAEP(d.generator, d.privateKey, msg)
}
