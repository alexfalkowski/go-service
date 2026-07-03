package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/message"
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
// key material or constructing the AES-GCM cipher is returned.
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

	aead, err := newAEAD(k)
	if err != nil {
		return nil, err
	}

	return &Cipher{gen: gen, aead: aead}, nil
}

// Cipher provides AES-GCM encryption and decryption using a configured key.
//
// The ciphertext format produced by Encrypt is:
//
//	nonce || gcm(ciphertext+tag)
//
// Where nonce is generated fresh per encryption and is required to decrypt.
type Cipher struct {
	gen  *rand.Generator
	aead cipher.AEAD
}

// Encrypt encrypts msg.Data using AES-GCM and returns nonce||ciphertext.
//
// A fresh nonce is generated for each call and is prefixed to the returned byte slice so Decrypt can recover it.
// The returned slice includes the GCM authentication tag as produced by [cipher.AEAD.Seal].
// msg.Meta is authenticated as AES-GCM associated data and must be supplied unchanged to Decrypt.
//
// Errors are returned if plaintext exceeds AES-GCM's per-nonce maximum or nonce generation fails.
func (c *Cipher) Encrypt(msg message.Message) ([]byte, error) {
	if int64(len(msg.Data)) > maxGCMPlaintextSize {
		return nil, ErrInvalidPlaintextLength
	}

	nonce, err := c.gen.GenerateBytes(c.aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return c.aead.Seal(nonce, nonce, msg.Data, msg.Meta), nil
}

// Decrypt decrypts a value produced by Encrypt.
//
// The msg.Data field must be in the format nonce||ciphertext, where nonce length is determined by the
// underlying AEAD (GCM) nonce size. msg.Meta must match the metadata supplied to Encrypt.
//
// Errors:
//   - ErrInvalidLength if msg.Data is shorter than the required nonce size.
//   - Any error returned by the underlying AEAD if authentication fails or if msg.Data is malformed.
func (c *Cipher) Decrypt(msg message.Message) ([]byte, error) {
	size := c.aead.NonceSize()
	if len(msg.Data) < size {
		return nil, ErrInvalidLength
	}

	nonce, text := msg.Data[:size], msg.Data[size:]
	return c.aead.Open(nil, nonce, text, msg.Meta)
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes: invalid key size %d: %w", len(key), crypto.ErrInvalidKeySize)
	}

	return cipher.NewGCM(b)
}
