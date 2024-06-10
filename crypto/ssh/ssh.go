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
func Generate() (PublicKey, PrivateKey, error) {
	pu, pr, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	ppr, err := ssh.MarshalPrivateKey(pr, "")
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(pu)
	if err != nil {
		return "", "", err
	}

	return PublicKey(ssh.MarshalAuthorizedKey(pub)), PrivateKey(pem.EncodeToMemory(ppr)), nil
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

func (a *sshAlgo) Sign(msg string) string {
	m := ed25519.Sign(a.privateKey, []byte(msg))

	return base64.StdEncoding.EncodeToString(m)
}

func (a *sshAlgo) Verify(sig, msg string) error {
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
