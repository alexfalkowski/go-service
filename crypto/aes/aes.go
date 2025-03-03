package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/alexfalkowski/go-service/crypto/rand"
)

// Code is adapted from https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865.

// ErrInvalidLength for aes.
var ErrInvalidLength = errors.New("aes: invalid length")

// NewGenerator for aes.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for aes.
type Generator struct {
	gen *rand.Generator
}

// Generate for aes.
func (g *Generator) Generate() (string, error) {
	return g.gen.GenerateText(32)
}

// NewCipher for aes.
func NewCipher(gen *rand.Generator, cfg *Config) (*Cipher, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	k, err := cfg.GetKey()

	return &Cipher{gen: gen, key: k}, err
}

// Cipher for aes.
type Cipher struct {
	gen *rand.Generator
	key []byte
}

// Encrypt for aes.
func (a *Cipher) Encrypt(msg string) (string, error) {
	aead, err := a.aead()
	if err != nil {
		return "", err
	}

	bytes, err := a.gen.GenerateBytes(aead.NonceSize())
	if err != nil {
		return "", err
	}

	s := aead.Seal(bytes, bytes, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(s), nil
}

// Decrypt for aes.
func (a *Cipher) Decrypt(msg string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	aead, err := a.aead()
	if err != nil {
		return "", err
	}

	size := aead.NonceSize()
	if len(decoded) < size {
		return "", ErrInvalidLength
	}

	nonce, c := decoded[:size], decoded[size:]
	decoded, err = aead.Open(nil, nonce, c, nil)

	return string(decoded), err
}

func (a *Cipher) aead() (cipher.AEAD, error) {
	b, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(b)
}
