package jwt

import (
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token/errors"
	"github.com/golang-jwt/jwt/v4"
)

// NewToken for jwt.
func NewToken(cfg *Config, signer *ed25519.Signer, verifier *ed25519.Verifier, generator id.Generator) *Token {
	if !IsEnabled(cfg) {
		return nil
	}

	return &Token{cfg: cfg, signer: signer, verifier: verifier, generator: generator}
}

// Token for jwt.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate JWT token.
func (t *Token) Generate() (string, error) {
	exp := time.MustParseDuration(t.cfg.Expiration)
	key := t.signer.PrivateKey
	now := time.Now()

	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: now.Add(exp)},
		ID:        t.generator.Generate(),
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    t.cfg.Issuer,
		NotBefore: &jwt.NumericDate{Time: now},
		Subject:   t.cfg.Subject,
		Audience:  []string{t.cfg.Audience},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = t.cfg.KeyID

	return token.SignedString(key)
}

// Verify JWT token.
func (t *Token) Verify(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, t.validate)
	if err != nil {
		return "", err
	}

	if !claims.VerifyIssuer(t.cfg.Issuer, true) {
		return "", errors.ErrInvalidIssuer
	}

	if !claims.VerifyAudience(t.cfg.Audience, true) {
		return "", errors.ErrInvalidAudience
	}

	return claims.Subject, nil
}

func (j *Token) validate(_ *jwt.Token) (any, error) {
	return j.verifier.PublicKey, nil
}
