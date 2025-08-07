package aes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// ErrInvalidLength for aes.
var ErrInvalidLength = errors.New("aes: invalid length")

// NewCipher for aes.
func NewCipher(gen *rand.Generator, fs *os.FS, cfg *Config) (*Cipher, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
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
