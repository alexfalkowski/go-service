package auth0

import (
	"errors"

	"github.com/form3tech-oss/jwt-go"
)

var (
	// ErrInvalidAudience for Auth0.
	ErrInvalidAudience = errors.New("invalid audience")

	// ErrInvalidIssuer for Auth0.
	ErrInvalidIssuer = errors.New("invalid issuer")

	// ErrInvalidAlgorithm for Auth0.
	ErrInvalidAlgorithm = errors.New("invalid algorithm")

	// ErrInvalidToken for Auth0.
	ErrInvalidToken = errors.New("invalid token")
)

type verifier struct {
	cfg  *Config
	cert Certificator
}

func (v *verifier) Verify(token []byte) error {
	parsedToken, err := jwt.Parse(string(token), v.validate)
	if err != nil {
		return err
	}

	if parsedToken.Header["alg"] != "RS256" {
		return ErrInvalidAlgorithm
	}

	if !parsedToken.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (v *verifier) validate(token *jwt.Token) (interface{}, error) {
	claims := token.Claims.(jwt.MapClaims)

	checkAud := claims.VerifyAudience(v.cfg.Audience, false)
	if !checkAud {
		return token, ErrInvalidAudience
	}

	checkIss := claims.VerifyIssuer(v.cfg.Issuer, false)
	if !checkIss {
		return token, ErrInvalidIssuer
	}

	cert, err := v.cert.Certificate(token)
	if err != nil {
		return token, err
	}

	result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return token, err
	}

	return result, nil
}
