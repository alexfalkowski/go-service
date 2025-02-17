package token

import (
	"encoding/hex"

	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/crypto/rand"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/time"
	"github.com/golang-jwt/jwt/v4"
)

// KID is a key ID.
type KID string

// NewKID for JWKSets.
func NewKID(gen *rand.Generator) (KID, error) {
	b, err := gen.GenerateLetters(10)
	if err != nil {
		return "", err
	}

	return KID(hex.EncodeToString([]byte(b))), nil
}

// JWT token.
type JWT struct {
	ed  *ed25519.Signer
	gen id.Generator
	kid KID
}

// NewJWT token.
func NewJWT(kid KID, ed *ed25519.Signer, gen id.Generator) *JWT {
	return &JWT{kid: kid, ed: ed, gen: gen}
}

// Generate JWT token.
func (j *JWT) Generate(sub, aud, iss string, exp time.Duration) (string, error) {
	key := j.ed.PrivateKey
	now := time.Now()

	claims := &jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: now.Add(exp)},
		ID:        j.gen.Generate(),
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
	return j.ed.PublicKey, nil
}
