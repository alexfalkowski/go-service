package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/context"
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
func (t *Token) Generate(ctx context.Context) (string, error) {
	opts := context.Opts(ctx)
	exp := time.MustParseDuration(t.cfg.Expiration)
	now := time.Now()

	token := paseto.NewToken()
	token.SetJti(t.generator.Generate())
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(exp))
	token.SetIssuer(t.cfg.Issuer)
	token.SetSubject(opts.GetString("sub"))
	token.SetAudience(opts.GetString("aud"))

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(t.signer.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.V4Sign(s, nil), nil
}

// Verify Paseto token.
func (t *Token) Verify(ctx context.Context, token string) (context.Context, error) {
	opts := context.Opts(ctx)
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(t.cfg.Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(opts.GetString("aud")))

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(t.verifier.PublicKey)
	if err != nil {
		return ctx, err
	}

	to, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return ctx, err
	}

	sub, err := to.GetSubject()
	if err != nil {
		return ctx, err
	}

	return context.AddToOpts(ctx, "sub", sub), nil
}
