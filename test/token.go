package test

import (
	"context"
	"errors"

	sjwt "github.com/alexfalkowski/go-service/pkg/security/jwt"
	"github.com/form3tech-oss/jwt-go"
)

// NewGenerator for test.
func NewGenerator(token string, err error) sjwt.Generator {
	return &generator{token: token, err: err}
}

type generator struct {
	token string
	err   error
}

func (g *generator) Generate(ctx context.Context) ([]byte, error) {
	return []byte(g.token), g.err
}

// NewVerifier for test.
func NewVerifier(token string) sjwt.Verifier {
	return &verifier{token: token}
}

type verifier struct {
	token string
}

// nolint:goerr113
func (v *verifier) Verify(ctx context.Context, token []byte) (*jwt.Token, error) {
	if string(token) != v.token {
		return nil, errors.New("invalid token")
	}

	jwtToken := &jwt.Token{
		Claims: jwt.MapClaims{
			"azp": v.token,
		},
	}

	return jwtToken, nil
}
