package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/time"
)

// TokenParams for paseto.
type TokenParams struct {
	di.In
	Config    *Config
	Signer    *ed25519.Signer
	Verifier  *ed25519.Verifier
	Generator id.Generator
}

// NewToken for paseto.
func NewToken(params TokenParams) *Token {
	if !params.Config.IsEnabled() {
		return nil
	}

	return &Token{
		cfg:       params.Config,
		signer:    params.Signer,
		verifier:  params.Verifier,
		generator: params.Generator,
	}
}

// Token for paseto.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate paseto token.
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
		return "", err
	}

	return token.V4Sign(s, nil), nil
}

// Verify Paseto token.
func (t *Token) Verify(token, aud string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(t.cfg.Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(aud))

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
