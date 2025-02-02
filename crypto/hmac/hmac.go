package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
)

// NewGenerator for hmac.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for hmac.
type Generator struct {
	gen *rand.Generator
}

// Generate for hmac.
func (g *Generator) Generate() (string, error) {
	s, err := g.gen.GenerateBytes(32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(s), nil
}

// NewSigner for hmac.
func NewSigner(cfg *Config) (Signer, error) {
	if !IsEnabled(cfg) {
		return &algo.NoSigner{}, nil
	}

	k, err := cfg.GetKey()

	return &signer{key: []byte(k)}, err
}

type (
	// Signer for hmac.
	Signer interface {
		algo.Signer
	}

	signer struct {
		key []byte
	}
)

func (a *signer) Sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, a.key)
	mac.Write([]byte(msg))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (a *signer) Verify(sig, msg string) error {
	decoded, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	mac := hmac.New(sha512.New, a.key)
	mac.Write([]byte(msg))

	expected := mac.Sum(nil)

	if !hmac.Equal(decoded, expected) {
		return errors.ErrInvalidMatch
	}

	return nil
}
