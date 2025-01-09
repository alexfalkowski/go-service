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
func Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	mpu, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return "", "", err
	}

	mpr, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return "", "", err
	}

	pub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu})
	pri := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr})

	return string(pub), string(pri), nil
}

// Algo for ed25519.
type Algo interface {
	algo.Signer

	// PublicKey for ed25519.
	PublicKey() ed25519.PublicKey

	// PrivateKey for ed25519.
	PrivateKey() ed25519.PrivateKey
}

// NewAlgo for ed25519.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &None{&algo.NoSigner{}}, nil
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

func (a *ed25519Algo) Sign(msg string) (string, error) {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m), nil
}

func (a *ed25519Algo) Verify(sig, msg string) error {
	d, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.publicKey, []byte(msg), d)
	if !ok {
		return errors.ErrInvalidMatch
	}

	return nil
}

// PublicKey for ed25519.
func (a *ed25519Algo) PublicKey() ed25519.PublicKey {
	return a.publicKey
}

// PrivateKey for ed25519.
func (a *ed25519Algo) PrivateKey() ed25519.PrivateKey {
	return a.privateKey
}

// None for ed25519.
type None struct {
	algo.Signer
}

// PublicKey for none.
func (a *None) PublicKey() ed25519.PublicKey {
	return nil
}

// PrivateKey for none.
func (a *None) PrivateKey() ed25519.PrivateKey {
	return nil
}
