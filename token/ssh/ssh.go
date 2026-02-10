package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type (
	// Signer aliases `crypto/ssh`.Signer.
	Signer = ssh.Signer

	// Verifier aliases `crypto/ssh`.Verifier.
	Verifier = ssh.Verifier
)

// NewToken constructs a Token using cfg and fs.
//
// If cfg is disabled, it returns nil.
func NewToken(cfg *Config, fs *os.FS) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, fs: fs}
}

// Token generates and verifies SSH-style tokens.
type Token struct {
	cfg *Config
	fs  *os.FS
}

// Generate creates an SSH-style token.
//
// The token format is `<name>-<base64(signature)>`, where signature is computed
// over the key name using the configured signing key.
func (t *Token) Generate() (string, error) {
	sig, err := ssh.NewSigner(t.fs, t.cfg.Key.Config)
	if err != nil {
		return strings.Empty, err
	}

	signature, err := sig.Sign(strings.Bytes(t.cfg.Key.Name))
	token := strings.Join("-", t.cfg.Key.Name, base64.Encode(signature))
	return token, err
}

// Verify validates token and returns the key name if it is valid.
//
// It expects token in the format `<name>-<base64(signature)>` and verifies the
// signature over `<name>` using the configured verification keys.
func (t *Token) Verify(token string) (string, error) {
	name, key, ok := strings.Cut(token, "-")
	if !ok {
		return strings.Empty, errors.ErrInvalidMatch
	}

	cfg := t.cfg.Keys.Get(name)
	if cfg == nil {
		return strings.Empty, errors.ErrInvalidMatch
	}

	verifier, err := ssh.NewVerifier(t.fs, cfg.Config)
	if err != nil {
		return strings.Empty, err
	}

	sig, err := base64.Decode(key)
	if err != nil {
		return strings.Empty, err
	}

	return name, verifier.Verify(sig, strings.Bytes(name))
}
