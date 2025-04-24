package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/errors"
)

// NewGenerator for ed25519.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator for ed25519.
type Generator struct {
	generator *rand.Generator
}

// Generate key pair with ed25519.
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.generator)
	if err != nil {
		return "", "", g.prefix(err)
	}

	mpu, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return "", "", g.prefix(err)
	}

	mpr, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return "", "", g.prefix(err)
	}

	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu}))
	pri := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr}))

	return pub, pri, nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("ed25519", err)
}
