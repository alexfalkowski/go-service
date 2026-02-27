package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
)

// NewToken constructs a Token that issues and validates PASETO v4 public (asymmetric) tokens.
//
// The resulting Token uses Ed25519 keys for signing and verification and an id.Generator
// for producing unique token IDs (jti). The keys are provided by the caller (typically
// via DI wiring).
//
// Enablement is modeled by presence: if cfg is nil, NewToken returns nil.
func NewToken(cfg *Config, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, signer: sig, verifier: ver, generator: gen}
}

// Token generates and verifies PASETO v4 public tokens.
//
// Issued tokens set common PASETO claims and identity fields (jti/iat/nbf/exp/iss/aud/sub).
//
// Note: This type assumes cfg, signer, verifier, and generator are non-nil. If you
// construct a Token with missing dependencies, methods may panic.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate creates a signed PASETO v4 public token for the given audience and subject.
//
// The token is signed using PASETO v4 public tokens (Ed25519 signatures). It sets
// common claims:
//
//   - jti: generated via the provided id.Generator
//   - iat: set to the current time
//   - nbf: set to the current time
//   - exp: set to now + parsed cfg.Expiration
//   - iss: from cfg.Issuer
//   - aud: set to the provided aud
//   - sub: set to the provided sub
//
// Expiration parsing uses time.MustParseDuration and will panic if cfg.Expiration is invalid.
// This is intended for fail-fast configuration behavior.
func (t *Token) Generate(aud, sub string) (string, error) {
	exp := time.MustParseDuration(t.cfg.Expiration)
	now := time.Now()
	token := paseto.NewToken()
	token.SetJti(t.generator.Generate())
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(exp))
	token.SetIssuer(t.cfg.Issuer)
	token.SetAudience(aud)
	token.SetSubject(sub)

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(t.signer.PrivateKey)
	if err != nil {
		return strings.Empty, err
	}

	return token.V4Sign(s, nil), nil
}

// Verify validates token and returns the subject (sub) if it is valid for the given audience.
//
// Verification is performed by constructing a parser with rules and then verifying the
// token signature using the configured Ed25519 public key.
//
// The rules enforced include:
//   - issuer matches cfg.Issuer (iss)
//   - token is not expired
//   - token is valid at the current time (iat/nbf semantics as defined by the upstream library)
//   - audience matches aud (aud)
//
// On failure, this method returns errors from the upstream PASETO library or from key
// construction. It does not currently map failures onto shared sentinel errors.
func (t *Token) Verify(token, aud string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(t.cfg.Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(aud))

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(t.verifier.PublicKey)
	if err != nil {
		return strings.Empty, err
	}

	to, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return strings.Empty, err
	}

	return to.GetSubject()
}
