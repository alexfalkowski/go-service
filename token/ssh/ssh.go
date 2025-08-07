package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type (
	// Signer is an alias of ssh.Signer.
	Signer = ssh.Signer

	// Verifier is an alias of ssh.Verifier.
	Verifier = ssh.Verifier
)

// NewToken for ssh.
func NewToken(fs *os.FS, cfg *Config) *Token {
	if !cfg.IsEnabled() {
		return nil
	}

	return &Token{fs: fs, cfg: cfg}
}

// Token for ssh.
type Token struct {
	fs  *os.FS
	cfg *Config
}

// Generate an SSH token.
func (t *Token) Generate() (string, error) {
	sig, err := ssh.NewSigner(t.fs, t.cfg.Key.Config)
	if err != nil {
		return "", err
	}

	signature, err := sig.Sign(strings.Bytes(t.cfg.Key.Name))
	token := strings.Join("-", t.cfg.Key.Name, base64.Encode(signature))
	return token, err
}

// Verify an SSH token.
func (t *Token) Verify(token string) (string, error) {
	name, key, ok := strings.Cut(token, "-")
	if !ok {
		return "", errors.ErrInvalidMatch
	}

	cfg := t.cfg.Keys.Get(name)
	if cfg == nil {
		return "", errors.ErrInvalidMatch
	}

	verifier, err := ssh.NewVerifier(t.fs, cfg.Config)
	if err != nil {
		return "", err
	}

	sig, err := base64.Decode(key)
	if err != nil {
		return "", err
	}

	return name, verifier.Verify(sig, strings.Bytes(name))
}
