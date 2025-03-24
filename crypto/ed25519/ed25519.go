package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	crypto "github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/errors"
)

// NewGenerator for ed25519.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for ed25519.
type Generator struct {
	gen *rand.Generator
}

// Generate key pair with ed25519.
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.gen)
	if err != nil {
		return "", "", errors.Prefix("ed25519", err)
	}

	mpu, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return "", "", errors.Prefix("ed25519", err)
	}

	mpr, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return "", "", errors.Prefix("ed25519", err)
	}

	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu}))
	pri := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr}))

	return pub, pri, nil
}

// NewSigner for ed25519.
func NewSigner(cfg *Config) (*Signer, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		return nil, err
	}

	pri, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}

	return &Signer{PublicKey: pub, PrivateKey: pri}, nil
}

// Signer for ed25519.
type Signer struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// Sign for ed25519.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}

// Verify for ed25519.
func (s *Signer) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(s.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
