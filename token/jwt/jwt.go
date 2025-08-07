package jwt

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/golang-jwt/jwt/v4"
)

// TokenParams for jwt.
type TokenParams struct {
	di.In
	Config    *Config
	Signer    *ed25519.Signer
	Verifier  *ed25519.Verifier
	Generator id.Generator
}

// NewToken for jwt.
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

// Token for jwt.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate JWT token.
func (t *Token) Generate(aud, sub string) (string, error) {
	exp := time.MustParseDuration(t.cfg.Expiration)
	key := t.signer.PrivateKey
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: now.Add(exp)},
		ID:        t.generator.Generate(),
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    t.cfg.Issuer,
		NotBefore: &jwt.NumericDate{Time: now},
		Audience:  []string{aud},
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = t.cfg.KeyID

	return token.SignedString(key)
}

// Verify JWT token.
func (t *Token) Verify(token, aud string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, t.validate)
	if err != nil {
		return "", err
	}

	if !claims.VerifyIssuer(t.cfg.Issuer, true) {
		return "", errors.ErrInvalidIssuer
	}

	if !claims.VerifyAudience(aud, true) {
		return "", errors.ErrInvalidAudience
	}

	return claims.Subject, nil
}

func (j *Token) validate(_ *jwt.Token) (any, error) {
	return j.verifier.PublicKey, nil
}
