package ssh

import (
	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/ssh"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	token "github.com/alexfalkowski/go-service/v2/token/errors"
)

const (
	tokenSeparator = "."
	tokenVersion   = "v1"
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
// This token kind is intentionally simple. It binds a logical key name and audience
// to a signature.
//
// Missing per-operation key material is treated as invalid configuration and
// reported via token/errors.ErrInvalidConfig.
type Token struct {
	cfg *Config
	fs  *os.FS
}

// Generate creates an SSH-style token for the given audience.
//
// Token format:
//
//	<base64(json-claims)>.<base64(signature)>
//
// Where json-claims contains:
//   - ver: the token format version ("v1")
//   - kid: t.cfg.Key.Name
//   - aud: aud
//   - iat: the issued-at Unix nanosecond timestamp
//   - exp: the expiration Unix nanosecond timestamp
//
// The sub parameter is accepted to match the other token implementations; SSH
// tokens identify the signing key instead of carrying a subject.
//
// The signature is produced by signing the exact JSON claims bytes using the
// configured signing key material. Because the key name and audience are both in
// the signed claims, a token minted for one audience cannot be replayed for
// another audience.
//
// High-level algorithm:
//  1. Load an SSH signer from the configured signing key (t.cfg.Key) using fs.
//  2. Marshal claims = {"ver": "v1", "kid": <name>, "aud": <aud>, "iat": <now>, "exp": <expiration>}.
//  3. Compute signature = Sign(claims).
//  4. Return "<base64(claims)>.<base64(signature)>".
//
// Errors are returned when the signing key configuration is missing/partial,
// key material cannot be loaded, claims encoding fails, or signature generation
// fails.
func (t *Token) Generate(aud, sub string) (string, error) {
	if t.cfg.Key == nil || t.cfg.Key.Config == nil || strings.IsEmpty(t.cfg.Key.Name) {
		return strings.Empty, token.ErrInvalidConfig
	}
	if t.cfg.Expiration <= 0 {
		return strings.Empty, token.ErrInvalidConfig
	}

	sig, err := ssh.NewSigner(t.fs, t.cfg.Key.Config)
	if err != nil {
		return strings.Empty, err
	}

	now := time.Now()
	c, err := json.Marshal(claims{
		Version:   tokenVersion,
		KeyID:     t.cfg.Key.Name,
		Audience:  aud,
		IssuedAt:  now.UnixNano(),
		ExpiresAt: now.Add(t.cfg.Expiration.Duration()).UnixNano(),
	})
	if err != nil {
		return strings.Empty, err
	}

	signature, err := sig.Sign(c)
	if err != nil {
		return strings.Empty, err
	}

	return strings.Join(tokenSeparator, base64.Encode(c), base64.Encode(signature)), nil
}

// Verify validates token for aud and returns the embedded key name if it is valid.
//
// Token format:
//
//	<base64(json-claims)>.<base64(signature)>
//
// Verification is name-based and audience-bound: Verify decodes the signed claims,
// checks claims.aud against aud, selects the matching verification key
// configuration from t.cfg.Keys using Keys.Get(claims.kid), and verifies the
// signature over the exact claims bytes with that key.
//
// High-level algorithm:
//  1. Split token into (encodedClaims, encodedSignature) on ".".
//  2. Decode and unmarshal the claims.
//  3. Look up a verification key config for claims.kid in t.cfg.Keys.
//  4. Load an SSH verifier from the selected key material using fs.
//  5. Decode the signature from base64.
//  6. Verify(signature, claims).
//  7. Check claims.ver, claims.aud, claims.iat, and claims.exp.
//
// Security-oriented error behavior:
//   - If the token cannot be split, the claims cannot be decoded, or no key
//     exists for the name, Verify returns crypto/errors.ErrInvalidMatch.
//     This intentionally collapses multiple invalid-token cases into a single
//     class to avoid leaking whether a given key name exists.
//   - If a matching key name exists but its verification config is missing/partial,
//     Verify returns token/errors.ErrInvalidConfig.
//   - Base64 decode errors and verifier loading errors are returned as-is.
//
// On success, Verify returns the extracted key name. On failure, it always
// returns an empty name alongside the error.
func (t *Token) Verify(tkn, aud string) (string, error) {
	data, encoded, key, err := parseClaims(tkn)
	if err != nil {
		return strings.Empty, err
	}

	cfg := t.cfg.Keys.Get(data.KeyID)
	if cfg == nil {
		return strings.Empty, crypto.ErrInvalidMatch
	}

	if cfg.Config == nil {
		return strings.Empty, token.ErrInvalidConfig
	}

	verifier, err := ssh.NewVerifier(t.fs, cfg.Config)
	if err != nil {
		return strings.Empty, err
	}

	sig, err := base64.Decode(key)
	if err != nil {
		return strings.Empty, err
	}

	if err := verifier.Verify(sig, encoded); err != nil {
		return strings.Empty, err
	}

	if err := validateClaims(data, aud, time.Now().UnixNano()); err != nil {
		return strings.Empty, err
	}

	return data.KeyID, nil
}
