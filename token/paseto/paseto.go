package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/time"
)

// NewToken for paseto.
func NewToken(cfg *Config, signer *ed25519.Signer, verifier *ed25519.Verifier, generator id.Generator) *Token {
	if !IsEnabled(cfg) {
		return nil
	}

	return &Token{cfg: cfg, signer: signer, verifier: verifier, generator: generator}
}

// Token for paseto.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate paseto token.
func (t *Token) Generate() (string, error) {
	exp := time.MustParseDuration(t.cfg.Expiration)
	now := time.Now()

	token := paseto.NewToken()
	token.SetJti(t.generator.Generate())
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(exp))
	token.SetIssuer(t.cfg.Issuer)
	token.SetSubject(t.cfg.Subject)
	token.SetAudience(t.cfg.Audience)

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(t.signer.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.V4Sign(s, nil), nil
}

// Verify Paseto token.
func (t *Token) Verify(token string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(t.cfg.Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(t.cfg.Audience))

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(t.verifier.PublicKey)
	if err != nil {
		return "", err
	}

	to, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return "", err
	}

	return to.GetSubject()
}
