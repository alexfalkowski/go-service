package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
	cerr "github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/runtime"
)

// NewGenerator for ed25519.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for hmac.
type Generator struct {
	gen *rand.Generator
}

// Generate key pair with Ed25519.
//
//nolint:nonamedreturns
func (g *Generator) Generate() (pub string, pri string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Prefix("ed25519", runtime.ConvertRecover(r))
		}
	}()

	public, private, err := ed25519.GenerateKey(g.gen)
	runtime.Must(err)

	mpu, err := x509.MarshalPKIXPublicKey(public)
	runtime.Must(err)

	mpr, err := x509.MarshalPKCS8PrivateKey(private)
	runtime.Must(err)

	pub = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: mpu}))
	pri = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mpr}))

	return
}

// NewSigner for ed25519.
func NewSigner(cfg *Config) (Signer, error) {
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

	return &signer{publicKey: pub, privateKey: pri}, nil
}

type (
	// Signer for ed25519.
	Signer interface {
		algo.Signer

		// PublicKey for ed25519.
		PublicKey() ed25519.PublicKey

		// PrivateKey for ed25519.
		PrivateKey() ed25519.PrivateKey
	}

	signer struct {
		publicKey  ed25519.PublicKey
		privateKey ed25519.PrivateKey
	}
)

func (a *signer) Sign(msg string) (string, error) {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m), nil
}

func (a *signer) Verify(sig, msg string) error {
	d, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.publicKey, []byte(msg), d)
	if !ok {
		return cerr.ErrInvalidMatch
	}

	return nil
}

// PublicKey for ed25519.
func (a *signer) PublicKey() ed25519.PublicKey {
	return a.publicKey
}

// PrivateKey for ed25519.
func (a *signer) PrivateKey() ed25519.PrivateKey {
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
