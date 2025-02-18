package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"

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
	return g.gen.GenerateText(32)
}

// NewSigner for hmac.
func NewSigner(cfg *Config) (*Signer, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	k, err := cfg.GetKey()

	return &Signer{key: []byte(k)}, err
}

// Signer for hmac.
type Signer struct {
	key []byte
}

// Sign for hmac.
func (a *Signer) Sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, a.key)
	mac.Write([]byte(msg))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Verify for hmac.
func (a *Signer) Verify(sig, msg string) error {
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
