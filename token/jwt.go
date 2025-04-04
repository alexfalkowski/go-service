package token

import (
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/time"
	"github.com/golang-jwt/jwt/v4"
)

const (
	// EmptyKID is the kid used when no kid is provided.
	EmptyKID = KID("")
)

// NewKID for JWKSets.
func NewKID(cfg *Config) KID {
	if !IsEnabled(cfg) {
		return EmptyKID
	}

	return KID(cfg.KeyID)
}

// KID is a key ID.
type KID string

// NewJWT token.
func NewJWT(kid KID, signer *ed25519.Signer, verifier *ed25519.Verifier, generator id.Generator) *JWT {
	return &JWT{kid: kid, signer: signer, verifier: verifier, generator: generator}
}

// JWT token.
type JWT struct {
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
	kid       KID
}

// Generate JWT token.
func (j *JWT) Generate(sub, aud, iss string, exp time.Duration) (string, error) {
	key := j.signer.PrivateKey
	now := time.Now()

	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: now.Add(exp)},
		ID:        j.generator.Generate(),
		IssuedAt:  &jwt.NumericDate{Time: now},
		Issuer:    iss,
		NotBefore: &jwt.NumericDate{Time: now},
		Subject:   sub,
		Audience:  []string{aud},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = j.kid

	return token.SignedString(key)
}

// Verify JWT token.
func (j *JWT) Verify(token, aud, iss string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, j.validate)
	if err != nil {
		return "", err
	}

	if !claims.VerifyIssuer(iss, true) {
		return "", ErrInvalidIssuer
	}

	if !claims.VerifyAudience(aud, true) {
		return "", ErrInvalidAudience
	}

	return claims.Subject, nil
}

func (j *JWT) validate(_ *jwt.Token) (any, error) {
	return j.verifier.PublicKey, nil
}
