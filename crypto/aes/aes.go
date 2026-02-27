package aes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// ErrInvalidLength is returned when a ciphertext is too short to contain the required nonce prefix.
var ErrInvalidLength = errors.New("aes: invalid length")

// NewCipher constructs an AES-GCM Cipher when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewCipher returns (nil, nil).
//
// Enabled behavior: the key material is loaded via cfg.GetKey(fs). Any error encountered while reading
// key material is returned.
//
// Note: this constructor does not validate the key length eagerly; key length validation occurs when
// Encrypt/Decrypt attempts to construct the underlying AES block cipher.
func NewCipher(gen *rand.Generator, fs *os.FS, cfg *Config) (*Cipher, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
	return &Cipher{gen: gen, key: k}, err
}

// Cipher provides AES-GCM encryption and decryption using a configured key.
//
// The ciphertext format produced by Encrypt is:
//
//	nonce || gcm(ciphertext+tag)
//
// Where nonce is generated fresh per encryption and is required to decrypt.
type Cipher struct {
	gen *rand.Generator
	key []byte
}

// Encrypt encrypts msg using AES-GCM and returns nonce||ciphertext.
//
// A fresh nonce is generated for each call and is prefixed to the returned byte slice so Decrypt can recover it.
// The returned slice includes the GCM authentication tag as produced by cipher.AEAD.Seal.
//
// Errors are returned if nonce generation fails or if the configured key is invalid for AES.
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
// The msg parameter must be in the format nonce||ciphertext, where nonce length is determined by the
// underlying AEAD (GCM) nonce size.
//
// Errors:
//   - ErrInvalidLength if msg is shorter than the required nonce size.
//   - Any error returned by the underlying AEAD if authentication fails or if msg is malformed.
//   - Any error returned if the configured key is invalid for AES.
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
