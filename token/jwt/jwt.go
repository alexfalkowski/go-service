package jwt

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/golang-jwt/jwt/v4"
)

// NewToken for jwt.
func NewToken(cfg *Config, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, signer: sig, verifier: ver, generator: gen}
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
		return strings.Empty, err
	}

	if !claims.VerifyIssuer(t.cfg.Issuer, true) {
		return strings.Empty, errors.ErrInvalidIssuer
	}

	if !claims.VerifyAudience(aud, true) {
		return strings.Empty, errors.ErrInvalidAudience
	}

	if err := claims.Valid(); err != nil {
		return strings.Empty, err
	}

	return claims.Subject, nil
}

func (j *Token) validate(token *jwt.Token) (any, error) {
	if token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
		return nil, errors.ErrInvalidAlgorithm
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || strings.IsEmpty(kid) {
		return nil, errors.ErrInvalidKeyID
	}

	if kid != j.cfg.KeyID {
		return nil, errors.ErrInvalidKeyID
	}

	return j.verifier.PublicKey, nil
}
