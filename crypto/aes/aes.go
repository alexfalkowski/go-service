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
var ErrInvalidLength = errors.New("invalid length")

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
	g, err := a.aead()
	if err != nil {
		return "", err
	}

	n, err := rand.GenerateBytes(uint32(g.NonceSize()))
	if err != nil {
		return "", err
	}

	s := g.Seal(n, n, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(s), nil
}

func (a *aesAlgo) Decrypt(msg string) (string, error) {
	d, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	g, err := a.aead()
	if err != nil {
		return "", err
	}

	size := g.NonceSize()
	if len(d) < size {
		return "", ErrInvalidLength
	}

	nonce, c := d[:size], d[size:]
	d, err = g.Open(nil, nonce, c, nil)

	return string(d), err
}

func (a *aesAlgo) aead() (cipher.AEAD, error) {
	b, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(b)
}
