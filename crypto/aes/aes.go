package aes

import (
	"crypto/aes"
	"crypto/cipher"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ErrInvalidLength is returned when a ciphertext is too short to contain the required nonce prefix.
var ErrInvalidLength = errors.New("aes: invalid length")

// ErrInvalidPlaintextLength is returned when plaintext exceeds AES-GCM's per-nonce maximum.
var ErrInvalidPlaintextLength = errors.New("aes: invalid plaintext length")

// maxGCMPlaintextSize is GCM's per-nonce plaintext limit: 2^32 - 2 AES blocks.
//
// GCM uses a 32-bit counter internally and reserves counter values, so the standard-library
// implementation panics above this limit. Check it before Seal so Encrypt can return an error.
const maxGCMPlaintextSize = ((1 << 32) - 2) * aes.BlockSize

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
	if strings.IsEmpty(cfg.Key) {
		return nil, crypto.ErrMissingKey
	}

	k, err := cfg.GetKey(fs)
	if err != nil {
		return nil, err
	}
	if len(k) == 0 {
		return nil, crypto.ErrMissingKey
	}

	return &Cipher{gen: gen, key: k}, nil
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
// Errors are returned if plaintext exceeds AES-GCM's per-nonce maximum, nonce generation fails,
// or the configured key is invalid for AES.
func (c *Cipher) Encrypt(msg []byte) ([]byte, error) {
	if int64(len(msg)) > maxGCMPlaintextSize {
		return nil, ErrInvalidPlaintextLength
	}

	aead, err := c.aead()
	if err != nil {
		return nil, err
	}

	nonce, err := c.gen.GenerateBytes(aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return aead.Seal(nonce, nonce, msg, nil), nil
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
