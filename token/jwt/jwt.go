package jwt

import (
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/golang-jwt/jwt/v4"
)

// NewToken constructs a Token that issues and validates JWTs according to cfg.
//
// The resulting Token uses Ed25519 keys for signing and verification and an id.Generator
// for producing unique JWT IDs (jti). The keys are provided by the caller (typically via DI).
//
// Enablement is modeled by presence: if cfg is nil, NewToken returns nil.
func NewToken(cfg *Config, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{cfg: cfg, signer: sig, verifier: ver, generator: gen}
}

// Token generates and verifies JWTs signed using Ed25519 (EdDSA).
//
// Issued tokens use standard registered claims (jwt.RegisteredClaims) and include a
// "kid" header to bind the token to a configured key identity.
//
// Note: This type assumes cfg, signer, verifier, and generator are non-nil. If you
// construct a Token with missing dependencies, methods may panic.
type Token struct {
	cfg       *Config
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate creates a signed JWT for the given audience and subject.
//
// The token is signed with Ed25519 using the JWT "EdDSA" signing method.
// It sets standard registered claims:
//
//   - iss: from cfg.Issuer
//   - aud: set to the provided aud (as a single-element audience list)
//   - sub: set to the provided sub
//   - iat: set to the current time
//   - nbf: set to the current time
//   - exp: set to now + parsed cfg.Expiration
//   - jti: generated via the provided id.Generator
//
// In addition, it sets the JWT header:
//
//   - kid: from cfg.KeyID
//
// Expiration parsing uses time.MustParseDuration and will panic if cfg.Expiration is invalid.
// This is intended for fail-fast configuration behavior.
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

// Verify validates token and returns the subject (sub) if it is valid for the given audience.
//
// Verification enforces the following checks:
//
//   - The token's signature algorithm is EdDSA (Ed25519).
//   - The JWT header "kid" exists, is non-empty, and matches cfg.KeyID exactly.
//   - The issuer claim ("iss") matches cfg.Issuer.
//   - The audience claim ("aud") contains the expected aud.
//   - Registered claim time validity using jwt.RegisteredClaims.Valid (exp/nbf/iat).
//
// This method returns sentinel errors from token/errors for some common classes of
// failures (issuer/audience mismatches and validate-time algorithm/kid mismatches).
// Parse/validation errors produced by the upstream JWT library may be returned as-is.
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
