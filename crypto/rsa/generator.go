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
// The generator parameter is retained for API consistency with other crypto
// generators. With Go 1.26 and later, crypto/rsa.GenerateKey uses the standard
// library's secure random source and ignores the supplied reader by default.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates RSA key pairs.
type Generator struct {
	generator *rand.Generator
}

// Generate returns an RSA public/private key pair encoded as PEM strings.
//
// The generated key size is KeySize bits.
// With Go 1.26 and later, crypto/rsa.GenerateKey uses the standard library's
// secure random source and ignores this Generator's injected reader by default.
//
// The returned PEM blocks are compatible with the expectations of [Config]:
//
//   - public:  a PEM block with Type "RSA PUBLIC KEY" containing PKCS#1-encoded bytes (x509.MarshalPKCS1PublicKey)
//   - private: a PEM block with Type "RSA PRIVATE KEY" containing PKCS#1-encoded bytes (x509.MarshalPKCS1PrivateKey)
//
// If key generation or encoding fails, the returned error is prefixed with "rsa".
func (g *Generator) Generate() (string, string, error) {
	privateKey, err := rsa.GenerateKey(g.generator, KeySize)
	if err != nil {
		return strings.Empty, strings.Empty, g.prefix(err)
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)})
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	return bytes.String(pub), bytes.String(pri), nil
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("rsa", err)
}
