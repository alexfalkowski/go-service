package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewGenerator for rsa.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator for rsa.
type Generator struct {
	generator *rand.Generator
}

// Generate key pair with rsa.
func (g *Generator) Generate() (string, string, error) {
	public, err := rsa.GenerateKey(g.generator, 4096)
	if err != nil {
		return strings.Empty, strings.Empty, err
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&public.PublicKey)})
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(public)})
	return bytes.String(pub), bytes.String(pri), nil
}
