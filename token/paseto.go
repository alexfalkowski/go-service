package token

import (
	"aidanwoods.dev/go-paseto"
	"github.com/alexfalkowski/go-service/crypto/ed25519"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/time"
)

// NewPaseto token.
func NewPaseto(signer *ed25519.Signer, verifier *ed25519.Verifier, generator id.Generator) *Paseto {
	return &Paseto{signer: signer, verifier: verifier, generator: generator}
}

// Paseto token.
type Paseto struct {
	signer    *ed25519.Signer
	verifier  *ed25519.Verifier
	generator id.Generator
}

// Generate Paseto token.
func (p *Paseto) Generate(sub, aud, iss string, exp time.Duration) (string, error) {
	now := time.Now()

	token := paseto.NewToken()
	token.SetJti(p.generator.Generate())
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(exp))
	token.SetIssuer(iss)
	token.SetSubject(sub)
	token.SetAudience(aud)

	s, err := paseto.NewV4AsymmetricSecretKeyFromBytes(p.signer.PrivateKey)
	if err != nil {
		return "", err
	}

	return token.V4Sign(s, nil), nil
}

// Verify Paseto token.
func (p *Paseto) Verify(token, aud, iss string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(iss))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	parser.AddRule(paseto.ForAudience(aud))

	s, err := paseto.NewV4AsymmetricPublicKeyFromBytes(p.verifier.PublicKey)
	if err != nil {
		return "", err
	}

	to, err := parser.ParseV4Public(s, token, nil)
	if err != nil {
		return "", err
	}

	return to.GetSubject()
}
