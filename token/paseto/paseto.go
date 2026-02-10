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
// If cfg is disabled (nil), it returns nil.
func NewToken(cfg *Config, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, signer: sig, verifier: ver, generator: gen}
}

// Token generates and verifies PASETO v4 public tokens.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate creates a signed PASETO v4 public token for the given audience and subject.
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

// Verify validates token and returns the subject if it is valid for the given audience.
//
// Verification enforces:
// - issuer matches cfg.Issuer
// - token is not expired and is valid at the current time
// - audience matches aud.
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
