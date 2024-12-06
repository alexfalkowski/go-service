package token

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/google/uuid"
)

// Paseto token.
type Paseto struct {
	ed ed25519.Algo
}

// NewPaseto token.
func NewPaseto(ed ed25519.Algo) *Paseto {
	return &Paseto{ed: ed}
}

// Generate Paseto token.
func (p *Paseto) Generate(sub, aud, iss string, exp time.Duration) (string, error) {
	k := p.ed.PrivateKey()
	t := time.Now()
	token := paseto.NewToken()

	token.SetJti(uuid.NewString())
	token.SetIssuedAt(t)
	token.SetNotBefore(t)
	token.SetExpiration(t.Add(exp))
	token.SetIssuer(iss)
	token.SetSubject(sub)
	token.SetAudience(aud)

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(k)
	if err != nil {
		return "", err
	}

	return token.V4Sign(s, nil), nil
}

// Verify Paseto token.
func (p *Paseto) Verify(token, aud, iss string) (string, error) {
	k := p.ed.PublicKey()

	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(iss))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(aud))

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(k)
	if err != nil {
		return "", err
	}

	to, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return "", err
	}

	return to.GetSubject()
}
