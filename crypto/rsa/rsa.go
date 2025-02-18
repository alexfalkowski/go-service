package rsa

import (
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/rand"
)

// NewGenerator for rsa.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for rsa.
type Generator struct {
	gen *rand.Generator
}

// Generate key pair with rsa.
func (g *Generator) Generate() (string, string, error) {
	public, err := rsa.GenerateKey(g.gen, 4096)
	if err != nil {
		return "", "", err
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&public.PublicKey)})
	pri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(public)})

	return string(pub), string(pri), nil
}

// NewCipher for rsa.
func NewCipher(gen *rand.Generator, cfg *Config) (*Cipher, error) {
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

	return &Cipher{gen: gen, publicKey: pub, privateKey: pri}, nil
}

// Cipher for rsa.
type Cipher struct {
	gen        *rand.Generator
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// Encrypt for rsa.
func (a *Cipher) Encrypt(msg string) (string, error) {
	e, err := rsa.EncryptOAEP(sha512.New(), a.gen, a.publicKey, []byte(msg), nil)

	return base64.StdEncoding.EncodeToString(e), err
}

// Decrypt for rsa.
func (a *Cipher) Decrypt(msg string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", err
	}

	decoded, err = rsa.DecryptOAEP(sha512.New(), a.gen, a.privateKey, decoded, nil)

	return string(decoded), err
}
