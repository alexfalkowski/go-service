package ssh

import (
	"crypto/ed25519"
	"encoding/pem"

	crypto "github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/errors"
	"golang.org/x/crypto/ssh"
)

// NewGenerator for ssh.
func NewGenerator(gen *rand.Generator) *Generator {
	return &Generator{gen: gen}
}

// Generator for ssh.
type Generator struct {
	gen *rand.Generator
}

// Generate key pair with ssh.
func (g *Generator) Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(g.gen)
	if err != nil {
		return "", "", errors.Prefix("ssh", err)
	}

	mpu, err := ssh.NewPublicKey(public)
	if err != nil {
		return "", "", errors.Prefix("ssh", err)
	}

	mpr, err := ssh.MarshalPrivateKey(private, "")
	if err != nil {
		return "", "", errors.Prefix("ssh", err)
	}

	pub := string(ssh.MarshalAuthorizedKey(mpu))
	pri := string(pem.EncodeToMemory(mpr))

	return pub, pri, nil
}

// NewSigner for ssh.
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

// Signer for ssh.
type Signer struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// Sign for ssh.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}

// Verify for ssh.
func (s *Signer) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(s.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
