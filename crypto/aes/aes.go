package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/rand"
)

// Code is adapted from https://gist.github.com/fracasula/38aa1a4e7481f9cedfa78a0cdd5f1865.

// ErrInvalidLength for aes.
var ErrInvalidLength = errors.New("aes: invalid length")

// Generate for aes.
func Generate() (string, error) {
	s, err := rand.GenerateBytes(32)

	return base64.StdEncoding.EncodeToString(s), err
}

// Algo for aes.
type Algo interface {
	algo.Cipher
}

// NewAlgo for aes.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &algo.NoCipher{}, nil
	}

	k, err := cfg.GetKey()

	return &aesAlgo{key: []byte(k)}, err
}

type aesAlgo struct {
	key []byte
}

func (a *aesAlgo) Encrypt(msg string) (string, error) {
	aead, err := a.aead()
	if err != nil {
		return "", err
	}

	bytes, err := rand.GenerateBytes(uint32(aead.NonceSize())) //nolint:gosec
	if err != nil {
		return "", err
	}

	s := aead.Seal(bytes, bytes, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(s), nil
}

func (a *aesAlgo) Decrypt(msg string) (string, error) {
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

func (a *aesAlgo) aead() (cipher.AEAD, error) {
	b, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(b)
}
