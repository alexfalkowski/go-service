package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
)

// Generate key pair with Ed25519.
func Generate() (PublicKey, PrivateKey, error) {
	pu, pr, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	mpu, err := x509.MarshalPKIXPublicKey(pu)
	if err != nil {
		return "", "", err
	}

	mpr, err := x509.MarshalPKCS8PrivateKey(pr)
	if err != nil {
		return "", "", err
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu})
	pri := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr})

	return PublicKey(pub), PrivateKey(pri), nil
}

// Algo for ed25519.
type Algo interface {
	algo.Signer
}

// NewAlgo for ed25519.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &algo.NoSigner{}, nil
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		return nil, err
	}

	pri, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}

	return &ed25519Algo{publicKey: pub, privateKey: pri}, nil
}

type ed25519Algo struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func (a *ed25519Algo) Sign(msg string) string {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m)
}

func (a *ed25519Algo) Verify(sig, msg string) error {
	d, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.publicKey, []byte(msg), d)
	if !ok {
		return errors.ErrMismatch
	}

	return nil
}
