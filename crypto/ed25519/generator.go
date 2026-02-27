package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewGenerator constructs a Generator that produces Ed25519 key pairs.
//
// The provided generator is used as the cryptographically-secure randomness source for key generation.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates Ed25519 key pairs.
type Generator struct {
	generator *rand.Generator
}

// Generate returns an Ed25519 public/private key pair encoded as PEM strings.
//
// The returned PEM blocks are compatible with the expectations of `crypto/ed25519.Config`:
//
//   - public:  a PEM block with Type "PUBLIC KEY" containing PKIX-encoded bytes (x509.MarshalPKIXPublicKey)
//   - private: a PEM block with Type "PRIVATE KEY" containing PKCS#8-encoded bytes (x509.MarshalPKCS8PrivateKey)
//
// If key generation, marshaling, or encoding fails, the returned error is prefixed with "ed25519".
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.generator)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	mpu, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	mpr, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu})
	pri := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr})
	return bytes.String(pub), bytes.String(pri), nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("ed25519", err)
}
