package aes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// ErrInvalidLength is returned when ciphertext is shorter than the nonce size.
var ErrInvalidLength = errors.New("aes: invalid length")

// NewCipher constructs an AES-GCM Cipher when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, the key material is loaded via cfg.GetKey.
func NewCipher(gen *rand.Generator, fs *os.FS, cfg *Config) (*Cipher, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
	return &Cipher{gen: gen, key: k}, err
}

// Cipher provides AES-GCM encryption and decryption using a configured key.
type Cipher struct {
	gen *rand.Generator
	key []byte
}

// Encrypt encrypts msg using AES-GCM and returns nonce||ciphertext.
//
// A fresh nonce is generated for each call and is prefixed to the returned byte slice so Decrypt can recover it.
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

// Decrypt decrypts a value produced by Encrypt.
//
// msg is expected to be nonce||ciphertext. If msg is shorter than the nonce size, it returns ErrInvalidLength.
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
