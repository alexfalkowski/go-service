package jwt

import (
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
	"github.com/golang-jwt/jwt/v4"
)

// Claims aliases the upstream JWT claims interface.
type Claims = jwt.Claims

// NumericDate aliases the upstream JWT numeric date type.
type NumericDate = jwt.NumericDate

// RegisteredClaims aliases the upstream JWT registered claims type.
type RegisteredClaims = jwt.RegisteredClaims

// SigningMethod aliases the upstream JWT signing method interface.
type SigningMethod = jwt.SigningMethod

// NewWithClaims aliases the upstream JWT token constructor.
var NewWithClaims = jwt.NewWithClaims

// SigningMethodEdDSA aliases the upstream JWT EdDSA signing method.
var SigningMethodEdDSA = jwt.SigningMethodEdDSA

// NewToken constructs a Token that issues and validates JWTs according to cfg.
//
// The resulting Token uses configured Ed25519 keys for signing and verification and an
// [id.Generator] for producing unique JWT IDs (jti).
//
// Enablement is modeled by presence: if cfg is nil, NewToken returns nil.
func NewToken(cfg *Config, fs *os.FS, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}

	return &Token{cfg: cfg, decoder: pem.NewDecoder(fs), generator: gen}
}

// Token generates and verifies JWTs signed using Ed25519 (EdDSA).
//
// Issued tokens use standard registered claims ([github.com/golang-jwt/jwt/v4.RegisteredClaims]) and include a
// "kid" header to bind the token to a configured key identity.
//
// Missing generation or verification dependencies are reported as [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
type Token struct {
	cfg       *Config
	decoder   *pem.Decoder
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
//   - jti: generated via the provided [id.Generator]
//
// In addition, it sets the JWT header:
//
//   - kid: from cfg.Key
func (t *Token) Generate(aud, sub string) (string, error) {
	if t.generator == nil {
		return strings.Empty, errors.ErrInvalidConfig
	}
	if strings.IsEmpty(t.cfg.Issuer) || strings.IsEmpty(t.cfg.Key) || t.cfg.Expiration <= 0 {
		return strings.Empty, errors.ErrInvalidConfig
	}

	key, err := t.cfg.Keys.Get(t.cfg.Key).Signer(t.decoder)
	if err != nil {
		return strings.Empty, err
	}

	now := time.Now()
	claims := &RegisteredClaims{
		ExpiresAt: &NumericDate{Time: now.Add(t.cfg.Expiration.Duration())},
		ID:        t.generator.Generate(),
		IssuedAt:  &NumericDate{Time: now},
		Issuer:    t.cfg.Issuer,
		NotBefore: &NumericDate{Time: now},
		Audience:  []string{aud},
		Subject:   sub,
	}
	token := NewWithClaims(SigningMethodEdDSA, claims)
	token.Header["kid"] = t.cfg.Key

	return token.SignedString(key.PrivateKey)
}

// Verify validates token and returns the subject (sub) if it is valid for the given audience.
//
// Verification enforces the following checks:
//
//   - The token's signature algorithm is EdDSA (Ed25519).
//   - The JWT header "kid" exists, is non-empty, and selects a configured verification key.
//   - The issuer claim ("iss") matches cfg.Issuer.
//   - The audience claim ("aud") contains the expected aud.
//   - Registered claim time validity with cfg.Leeway clock-skew tolerance (exp/nbf/iat).
//   - The signed lifetime (exp - iat) does not exceed cfg.Expiration.
//
// This method returns sentinel errors from token/errors for some common classes of
// failures (issuer/audience mismatches and validate-time algorithm/kid mismatches).
// Parse/validation errors produced by the upstream JWT library may be returned as-is.
// Upstream parse and signature errors can be returned before the local issuer,
// audience, required-claim, time-validity, and lifetime checks run.
func (t *Token) Verify(token, aud string) (string, error) {
	if strings.IsEmpty(t.cfg.Issuer) || t.cfg.Expiration <= 0 {
		return strings.Empty, errors.ErrInvalidConfig
	}

	claims := &RegisteredClaims{}

	_, err := jwt.ParseWithClaims(token, claims, t.validate, jwt.WithoutClaimsValidation())
	if err != nil {
		return strings.Empty, err
	}

	if !claims.VerifyIssuer(t.cfg.Issuer, true) {
		return strings.Empty, errors.ErrInvalidIssuer
	}

	if !claims.VerifyAudience(aud, true) {
		return strings.Empty, errors.ErrInvalidAudience
	}

	if err := validateRequiredClaims(claims); err != nil {
		return strings.Empty, err
	}
	if err := validateTime(claims, t.cfg.Expiration, t.cfg.Leeway); err != nil {
		return strings.Empty, err
	}

	return claims.Subject, nil
}

func (t *Token) validate(token *jwt.Token) (any, error) {
	if token.Method.Alg() != SigningMethodEdDSA.Alg() {
		return nil, errors.ErrInvalidAlgorithm
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || strings.IsEmpty(kid) {
		return nil, errors.ErrInvalidKeyID
	}

	key := t.cfg.Keys.Get(kid)
	if key == nil {
		return nil, errors.ErrInvalidKeyID
	}

	verifier, err := key.Verifier(t.decoder)
	if err != nil {
		return nil, err
	}

	return verifier.PublicKey, nil
}

func validateRequiredClaims(claims *RegisteredClaims) error {
	if claims.ExpiresAt == nil || claims.IssuedAt == nil || claims.NotBefore == nil {
		return errors.ErrInvalidTime
	}

	if strings.IsEmpty(claims.Subject) {
		return errors.ErrInvalidSubject
	}

	return nil
}

func validateTime(claims *RegisteredClaims, maxLifetime, leeway time.Duration) error {
	now := time.Now()
	allowedFuture := now.Add(leeway.Duration())
	if claims.IssuedAt.After(allowedFuture) || claims.NotBefore.After(allowedFuture) {
		return errors.ErrInvalidTime
	}

	if !claims.ExpiresAt.Time.Add(leeway.Duration()).After(now) {
		return errors.ErrInvalidTime
	}

	return validateLifetime(claims, maxLifetime)
}

func validateLifetime(claims *RegisteredClaims, maxLifetime time.Duration) error {
	if !claims.ExpiresAt.After(claims.IssuedAt.Time) || claims.ExpiresAt.Sub(claims.IssuedAt.Time) > maxLifetime.Duration() {
		return errors.ErrInvalidTime
	}

	return nil
}
