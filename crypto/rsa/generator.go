package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewGenerator constructs a Generator that produces RSA key pairs.
//
// The provided generator is used as the cryptographically-secure randomness source for key generation.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates RSA key pairs.
type Generator struct {
	generator *rand.Generator
}

// Generate returns an RSA public/private key pair encoded as PEM strings.
//
// The generated key size is 4096 bits.
//
// The returned PEM blocks are compatible with the expectations of `crypto/rsa.Config`:
//
//   - public:  a PEM block with Type "RSA PUBLIC KEY" containing PKCS#1-encoded bytes (x509.MarshalPKCS1PublicKey)
//   - private: a PEM block with Type "RSA PRIVATE KEY" containing PKCS#1-encoded bytes (x509.MarshalPKCS1PrivateKey)
//
// If key generation or encoding fails, the returned error is prefixed with "rsa".
func (g *Generator) Generate() (string, string, error) {
	public, err := rsa.GenerateKey(g.generator, 4096)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&public.PublicKey)})
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(public)})
	return bytes.String(pub), bytes.String(pri), nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("rsa", err)
}
