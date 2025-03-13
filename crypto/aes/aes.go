package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/alexfalkowski/go-service/crypto/rand"
)

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
func (c *Cipher) Encrypt(msg []byte) ([]byte, error) {
	aead, err := c.aead()
	if err != nil {
		return nil, err
	}

	bytes, err := c.gen.GenerateBytes(aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return aead.Seal(bytes, bytes, msg, nil), nil
}

// Decrypt for aes.
func (c *Cipher) Decrypt(msg []byte) ([]byte, error) {
	aead, err := c.aead()
	if err != nil {
		return nil, err
	}

	size := aead.NonceSize()
	if len(msg) < size {
		return nil, ErrInvalidLength
	}

	nonce, text := msg[:size], msg[size:]

	return aead.Open(nil, nonce, text, nil)
}

func (c *Cipher) aead() (cipher.AEAD, error) {
	b, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(b)
}
