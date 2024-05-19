package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/alexfalkowski/go-service/crypto/errors"
)

// Generate key pair with Ed25519.
func Generate() (PublicKey, PrivateKey, error) {
	pub, pri, err := ed25519.GenerateKey(rand.Reader)

	return PublicKey(base64.StdEncoding.EncodeToString(pub)), PrivateKey(base64.StdEncoding.EncodeToString(pri)), err
}

// Algo for ed25519.
type Algo interface {
	// Generate an encoded msg.
	Generate(msg string) string

	// Compare encoded with msg.
	Compare(enc, msg string) error

	// PublicKey for Ed25519.
	PublicKey() ed25519.PublicKey

	// PrivateKey for Ed25519.
	PrivateKey() ed25519.PrivateKey
}

// NewAlgo for ed25519.
func NewAlgo(cfg *Config) (Algo, error) {
	if !IsEnabled(cfg) {
		return &none{}, nil
	}

	pub, err := cfg.GetPublic()
	if err != nil {
		return nil, err
	}

	pri, err := cfg.GetPrivate()
	if err != nil {
		return nil, err
	}

	return &algo{publicKey: pub, privateKey: pri}, nil
}

type algo struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func (a *algo) Generate(msg string) string {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m)
}

func (a *algo) Compare(enc, msg string) error {
	d, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return err
	}

	ok := ed25519.Verify(a.publicKey, []byte(msg), d)
	if !ok {
		return errors.ErrMismatch
	}

	return nil
}

func (a *algo) PublicKey() ed25519.PublicKey {
	return a.publicKey
}

func (a *algo) PrivateKey() ed25519.PrivateKey {
	return a.privateKey
}

type none struct{}

func (*none) Generate(msg string) string {
	return msg
}

func (*none) Compare(_, _ string) error {
	return nil
}

func (*none) PublicKey() ed25519.PublicKey {
	return ed25519.PublicKey{}
}

func (*none) PrivateKey() ed25519.PrivateKey {
	return ed25519.PrivateKey{}
}
