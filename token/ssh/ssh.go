package ssh

import (
	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Signer is an alias for crypto/ssh.Signer.
//
// It represents an object capable of producing signatures using SSH key material.
type Signer = ssh.Signer

// Verifier is an alias for crypto/ssh.Verifier.
//
// It represents an object capable of verifying signatures produced by an SSH signer.
type Verifier = ssh.Verifier

// NewToken constructs a Token using cfg and fs.
//
// The returned Token loads key material using fs when generating and verifying tokens.
//
// Enablement is modeled by configuration: if cfg is disabled (see Config.IsEnabled),
// NewToken returns nil.
func NewToken(cfg *Config, fs *os.FS) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, fs: fs}
}

// Token generates and verifies SSH-style tokens.
//
// This token kind is intentionally simple and does not carry claims (audience, issuer,
// expiration, etc.). Instead, it binds a logical key name to a signature.
//
// Note: Token assumes cfg and fs are non-nil and that cfg contains the appropriate
// key material for the operation (signing key for Generate, verification keys for
// Verify). If those dependencies are missing, methods may panic.
type Token struct {
	cfg *Config
	fs  *os.FS
}

// Generate creates an SSH-style token.
//
// Token format:
//
//	<name>-<base64(signature)>
//
// Where <name> is t.cfg.Key.Name and signature is produced by signing the bytes of
// <name> using the configured signing key material.
//
// High-level algorithm:
//  1. Load an SSH signer from the configured signing key (t.cfg.Key) using fs.
//  2. Compute signature = Sign(<name>).
//  3. Return "<name>-<base64(signature)>".
//
// Errors are returned when key material cannot be loaded or signature generation fails.
func (t *Token) Generate() (string, error) {
	sig, err := ssh.NewSigner(t.fs, t.cfg.Key.Config)
	if err != nil {
		return strings.Empty, err
	}

	signature, err := sig.Sign(strings.Bytes(t.cfg.Key.Name))
	token := strings.Join("-", t.cfg.Key.Name, base64.Encode(signature))
	return token, err
}

// Verify validates token and returns the embedded key name if it is valid.
//
// Token format:
//
//	<name>-<base64(signature)>
//
// Verification is name-based: Verify extracts <name> and then selects the matching
// verification key configuration from t.cfg.Keys using Keys.Get(name). It verifies
// the signature over the bytes of <name> with that key.
//
// High-level algorithm:
//  1. Split token into (name, encodedSignature) on the first "-".
//  2. Look up a verification key config for name in t.cfg.Keys.
//  3. Load an SSH verifier from the selected key material using fs.
//  4. Decode the signature from base64.
//  5. Verify(signature, <name>).
//
// Security-oriented error behavior:
//   - If the token cannot be split or no key exists for the name, Verify returns
//     errors.ErrInvalidMatch. This intentionally collapses multiple invalid-token
//     cases into a single class to avoid leaking whether a given key name exists.
//   - Base64 decode errors and verifier loading errors are returned as-is.
//
// On success, Verify returns the extracted name.
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
