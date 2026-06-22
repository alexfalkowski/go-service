package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token/errors"
)

// NewToken constructs a Token that issues and validates PASETO v4 public (asymmetric) tokens.
//
// The resulting Token uses configured Ed25519 keys for signing and verification and an
// [id.Generator] for producing unique token IDs (jti).
//
// Enablement is modeled by presence: if cfg is nil, NewToken returns nil.
func NewToken(cfg *Config, fs *os.FS, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}

	return &Token{cfg: cfg, decoder: pem.NewDecoder(fs), generator: gen}
}

// Token generates and verifies PASETO v4 public tokens.
//
// Issued tokens set common PASETO claims and identity fields (jti/iat/nbf/exp/iss/aud/sub).
//
// Missing generation or verification dependencies are reported as [github.com/alexfalkowski/go-service/v2/token/errors.ErrInvalidConfig].
type Token struct {
	cfg       *Config
	decoder   *pem.Decoder
	generator id.Generator
}

// Generate creates a signed PASETO v4 public token for the given audience and subject.
//
// The token is signed using PASETO v4 public tokens (Ed25519 signatures). It sets
// common claims:
//
//   - jti: generated via the provided [id.Generator]
//   - iat: set to the current time
//   - nbf: set to the current time
//   - exp: set to now + parsed cfg.Expiration
//   - iss: from cfg.Issuer
//   - aud: set to the provided aud
//   - sub: set to the provided sub
//   - footer kid: from cfg.Key
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
	token := paseto.NewToken()
	token.SetJti(t.generator.Generate())
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(t.cfg.Expiration.Duration()))
	token.SetIssuer(t.cfg.Issuer)
	token.SetAudience(aud)
	token.SetSubject(sub)

	rawFooter, err := encodeFooter(&footer{KeyID: t.cfg.Key})
	if err != nil {
		return strings.Empty, err
	}
	token.SetFooter(rawFooter)

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(key.PrivateKey)
	if err != nil {
		return strings.Empty, err
	}

	return token.V4Sign(s, nil), nil
}

// Verify validates token and returns the subject (sub) if it is valid for the given audience.
//
// Verification is performed by constructing a parser with rules and then verifying the
// token signature using the configured Ed25519 public key.
//
// The rules enforced include:
//   - issuer matches cfg.Issuer (iss)
//   - token is valid at the current time with cfg.Leeway clock-skew tolerance (iat/nbf/exp)
//   - audience matches aud (aud)
//   - the signed lifetime (exp - iat) does not exceed cfg.Expiration
//
// On failure, parser, rule, signature, and key-construction errors may come from
// the upstream PASETO library. Local config, subject, and signed-lifetime checks
// return shared sentinel errors from token/errors.
func (t *Token) Verify(token, aud string) (string, error) {
	if strings.IsEmpty(t.cfg.Issuer) || t.cfg.Expiration <= 0 {
		return strings.Empty, errors.ErrInvalidConfig
	}

	parser := paseto.NewParserWithoutExpiryCheck()
	parser.AddRule(paseto.IssuedBy(t.cfg.Issuer))
	parser.AddRule(paseto.ForAudience(aud))

	key, err := t.publicKey(token, parser)
	if err != nil {
		return strings.Empty, err
	}

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(key)
	if err != nil {
		return strings.Empty, err
	}

	parsed, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return strings.Empty, err
	}

	if err := validateTime(parsed, t.cfg.Expiration, t.cfg.Leeway); err != nil {
		return strings.Empty, err
	}

	return subject(parsed)
}

func (t *Token) publicKey(token string, parser paseto.Parser) ([]byte, error) {
	raw, err := parser.UnsafeParseFooter(paseto.V4Public, token)
	if err != nil {
		return nil, err
	}

	footer, err := parseFooter(raw)
	if err != nil {
		return nil, err
	}

	key := t.cfg.Keys.Get(footer.KeyID)
	if key == nil {
		return nil, errors.ErrInvalidKeyID
	}

	verifier, err := key.Verifier(t.decoder)
	if err != nil {
		return nil, err
	}

	return verifier.PublicKey, nil
}

func subject(token *paseto.Token) (string, error) {
	sub, err := token.GetSubject()
	if err != nil {
		return strings.Empty, err
	}
	if strings.IsEmpty(sub) {
		return strings.Empty, errors.ErrInvalidSubject
	}

	return sub, nil
}

func validateTime(token *paseto.Token, maxLifetime, leeway time.Duration) error {
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return errors.ErrInvalidTime
	}
	notBefore, err := token.GetNotBefore()
	if err != nil {
		return errors.ErrInvalidTime
	}
	expiresAt, err := token.GetExpiration()
	if err != nil {
		return errors.ErrInvalidTime
	}

	now := time.Now()
	allowedFuture := now.Add(leeway.Duration())
	if issuedAt.After(allowedFuture) || notBefore.After(allowedFuture) {
		return errors.ErrInvalidTime
	}
	if !expiresAt.Add(leeway.Duration()).After(now) {
		return errors.ErrInvalidTime
	}

	return validateLifetimeRange(issuedAt, expiresAt, maxLifetime)
}

func validateLifetimeRange(issuedAt, expiresAt time.Time, maxLifetime time.Duration) error {
	if !expiresAt.After(issuedAt) || expiresAt.Sub(issuedAt) > maxLifetime.Duration() {
		return errors.ErrInvalidTime
	}

	return nil
}
