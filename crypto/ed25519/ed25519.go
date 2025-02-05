package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

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

func (a *Signer) Sign(msg string) (string, error) {
	m := ed25519.Sign(a.PrivateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m), nil
}

func (a *Signer) Verify(sig, msg string) error {
	d, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.PublicKey, []byte(msg), d)
	if !ok {
		return cerr.ErrInvalidMatch
	}

	return nil
}
