package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"

	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"golang.org/x/crypto/ssh"
)

// Generate key pair with ssh.
func Generate() (string, string, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	pri, err := ssh.MarshalPrivateKey(private, "")
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(public)
	if err != nil {
		return "", "", err
	}

	return string(ssh.MarshalAuthorizedKey(pub)), string(pem.EncodeToMemory(pri)), nil
}

// Algo for ssh.
type Algo interface {
	algo.Signer
}

// NewAlgo for ssh.
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

	return &sshAlgo{publicKey: pub, privateKey: pri}, nil
}

type sshAlgo struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func (a *sshAlgo) Sign(msg string) (string, error) {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m), nil
}

func (a *sshAlgo) Verify(sig, msg string) error {
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
