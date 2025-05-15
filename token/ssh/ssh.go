package ssh

import (
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/alexfalkowski/go-service/crypto/ssh"
	"github.com/alexfalkowski/go-service/encoding/base64"
	"github.com/alexfalkowski/go-service/strings"
)

// NewToken for ssh.
func NewToken(cfg *Config) *Token {
	if !IsEnabled(cfg) {
		return nil
	}

	return &Token{cfg: cfg}
}

// Token for ssh.
type Token struct {
	cfg *Config
}

// Generate an SSH token.
func (t *Token) Generate() (string, error) {
	sig, err := ssh.NewSigner(t.cfg.Key.Config)
	if err != nil {
		return "", err
	}

	signature, err := sig.Sign(strings.Bytes(t.cfg.Key.Name))
	token := strings.Join("-", t.cfg.Key.Name, base64.Encode(signature))

	return token, err
}

// Verify an SSH token.
func (t *Token) Verify(token string) error {
	name, key, ok := strings.Cut(token, "-")
	if !ok {
		return errors.ErrInvalidMatch
	}

	cfg := t.cfg.Keys.Get(name)
	if cfg == nil {
		return errors.ErrInvalidMatch
	}

	verifier, err := ssh.NewVerifier(cfg.Config)
	if err != nil {
		return err
	}

	sig, err := base64.Decode(key)
	if err != nil {
		return err
	}

	return verifier.Verify(sig, strings.Bytes(name))
}
